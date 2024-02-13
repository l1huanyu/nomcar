package query

import "context"

type GetCarListHandler struct {
}

type GetCarListCmd struct {
	OwnerOpenID string
}

func (h *GetCarListHandler) Handle(ctx context.Context, cmd *GetCarListCmd) ([]string, error) {
	return nil, nil
}
