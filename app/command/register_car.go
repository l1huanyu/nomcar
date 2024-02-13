package command

import (
	"context"
)

type RegisterCarHandler struct {
}

type RegisterCarCmd struct {
	CarID         string
	OwnerOpenID   string
	OwnerPhoneNum int64
}

func (r RegisterCarHandler) Handle(ctx context.Context, cmd *RegisterCarCmd) error {
	return nil
}
