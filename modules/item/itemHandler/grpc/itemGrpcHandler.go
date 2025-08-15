package itemHandler

import (
	"context"
	itemPb "microService/modules/item/itemPb"
	"microService/modules/item/itemUsecase"
)

type (
	itemGrpcHandler struct {
		itemPb.UnimplementedItemGrpcServiceServer
		itemUsecase itemUsecase.ItemUsecase
	}
)

func NewitemGrpcHandler(itemUsecase itemUsecase.ItemUsecase) *itemGrpcHandler {
	return &itemGrpcHandler{itemUsecase: itemUsecase}
}

func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return nil, nil
}
