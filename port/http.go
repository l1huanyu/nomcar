package port

import (
	"fmt"
	"net/http"

	"github.com/l1huanyu/nomcar/app"
	"github.com/l1huanyu/nomcar/app/command"
	"github.com/l1huanyu/nomcar/app/query"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const SuccessHTML = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
	body {
		background-color: #f8f9fa;
		padding: 20px;
		text-align: center;
	}
	h1 {
		color: #343a40;
	}
	p {
		color: #6c757d;
	}
	</style>
	<title>「小🍐挪车」成功界面</title>
</head>
<body>
	<h1>成功！</h1>
	<p>「小🍐挪车」已为您%s通知车主挪车。</p>
</body>
</html>`

type HTTPHandler struct {
	a *app.App
	e *echo.Echo
}

func NewHTTPHandler() HTTPHandler {
	h := HTTPHandler{
		a: app.NewApp(),
		e: echo.New(),
	}
	h.registerRouter()
	return h
}

func (h HTTPHandler) registerRouter() {
	h.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)),
		middleware.Logger(), middleware.CORS(), middleware.Recover(), middleware.RequestID())
	// 注册
	h.e.POST("/nomcar/api/register", h.HandleRegisterCar)
	// 查看已注册列表
	h.e.GET("/nomcar/api/cars", h.HandleGetCarList)
	// 查看二维码
	h.e.GET("/nomcar/api/qrcode/:car", h.HandleGetCarQRCode)
	// 通知车主
	h.e.GET("/nomcar/api/notify/:car", h.HandleNotifyCarOwner)
}

func (h HTTPHandler) Run() {
	h.e.Logger.Fatal(h.e.Start(":8823"))
}

type RegisterCarReq struct {
	CarID         string `json:"car_id"`
	OwnerPhoneNum int64  `json:"owner_phone_num"`
}

// HandleRegisterCar 小程序端调用，注册.
func (h HTTPHandler) HandleRegisterCar(ctx echo.Context) error {
	req := new(RegisterCarReq)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "[HTTPHandler::HandleRegisterCar]ctx.Bind error,err="+err.Error())
	}
	ownerOpenID, err := parseOwnerOpenID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	cmd := command.RegisterCarCmd{
		CarID:         req.CarID,
		OwnerOpenID:   ownerOpenID,
		OwnerPhoneNum: req.OwnerPhoneNum,
	}
	if err := h.a.Commands.RegisterCar.Handle(ctx.Request().Context(), &cmd); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	logrus.Infof("[HTTPHandler::HandleRegisterCar]success,cmd=%+v", cmd)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// HandleGetCarList 小程序端调用，获取🚗列表.
func (h HTTPHandler) HandleGetCarList(ctx echo.Context) error {
	ownerOpenID, err := parseOwnerOpenID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	cmd := query.GetCarListCmd{
		OwnerOpenID: ownerOpenID,
	}
	cars, err := h.a.Queries.GetCarList.Handle(ctx.Request().Context(), &cmd)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
	}
	logrus.Infof("[HTTPHandler::HandleGetCarList]success,cmd=%+v", cmd)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"cars": cars,
		},
	})
}

// HandleGetCarQRCode 小程序端调用，生成端外二维码.
func (h HTTPHandler) HandleGetCarQRCode(ctx echo.Context) error {
	carID := ctx.Param("car")
	if carID == "" {
		return ctx.JSON(http.StatusBadRequest, "bad request, no car id")
	}
	cmd := query.GetCarQRCodeCmd{
		CarID: carID,
	}
	qrCode, err := h.a.Queries.GetCarQRCode.Handle(ctx.Request().Context(), &cmd)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
	}
	logrus.Infof("[HTTPHandler::HandleGetCarQRCode]success,cmd=%+v", cmd)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"qr_code": qrCode,
		},
	})
}

// HandleNotifyCarOwner 端外调用，发送小程序端内订阅消息.
func (h HTTPHandler) HandleNotifyCarOwner(ctx echo.Context) error {
	carID := ctx.Param("car")
	if carID == "" {
		return ctx.JSON(http.StatusBadRequest, "bad request, no car id")
	}
	channel := ctx.QueryParam("channel")
	if channel == "" {
		return ctx.JSON(http.StatusBadRequest, "bad request, no channel")
	}
	cmd := command.NotifyCarOwnerCmd{
		CarID:   carID,
		Channel: channel,
	}
	if err := h.a.Commands.NotifyCar.Handle(ctx.Request().Context(), &cmd); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
	}
	logrus.Infof("[HTTPHandler::HandleNotifyCarOwner]success,cmd=%+v", cmd)
	return ctx.HTML(http.StatusOK, fmt.Sprintf(SuccessHTML, channel))
}

func parseOwnerOpenID(ctx echo.Context) (string, error) {
	if ctx.Request() == nil || ctx.Request().Header == nil {
		return "", errors.New("bad request, no request header")
	}
	ownerOpenID := ctx.Request().Header.Get("x-wx-openid")
	if ownerOpenID == "" {
		return "", errors.New("bad request, no owner openid")
	}
	return ownerOpenID, nil
}
