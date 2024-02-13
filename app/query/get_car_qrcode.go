package query

import "context"

type GetCarQRCodeHandler struct {
}

type GetCarQRCodeCmd struct {
	CarID string
}

func (g GetCarQRCodeHandler) Handle(ctx context.Context, cmd *GetCarQRCodeCmd) (string, error) {
	return "", nil
}
