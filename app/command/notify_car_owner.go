package command

import "context"

type NotifyCarOwnerHandler struct {
}

type NotifyCarOwnerCmd struct {
	CarID   string
	Channel string
}

func (n NotifyCarOwnerHandler) Handle(ctx context.Context, cmd *NotifyCarOwnerCmd) error {
	return nil
}
