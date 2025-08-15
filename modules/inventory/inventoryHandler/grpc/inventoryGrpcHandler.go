package inventoryHandler

import (
	"context"
	inventoryPb "microService/modules/inventory/inventoryPb"
	"microService/modules/inventory/inventoryUsecase"
)

type (
	inventoryGrpcHandler struct {
		inventoryPb.UnimplementedInventoryGrpcServiceServer
		inventoryUsecase inventoryUsecase.InventoryUsecase
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecase) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUsecase: inventoryUsecase}
}

func (g *inventoryGrpcHandler) IsAvailableToSell(ctx context.Context, req *inventoryPb.IsAvailableToSellReq) (*inventoryPb.IsAvailableToSellRes, error) {
	return nil, nil
}
