package app

import (
	"github.com/l1huanyu/nomcar/app/command"
	"github.com/l1huanyu/nomcar/app/query"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterCar command.RegisterCarHandler
	NotifyCar   command.NotifyCarOwnerHandler
}

type Queries struct {
	GetCarQRCode query.GetCarQRCodeHandler
	GetCarList   query.GetCarListHandler
}

func NewApp() *App {
	return &App{
		Commands: Commands{
			RegisterCar: command.RegisterCarHandler{},
			NotifyCar:   command.NotifyCarOwnerHandler{},
		},
		Queries: Queries{
			GetCarQRCode: query.GetCarQRCodeHandler{},
			GetCarList:   query.GetCarListHandler{},
		},
	}
}
