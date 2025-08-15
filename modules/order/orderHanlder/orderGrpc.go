package orderhanlder

import (
	"context"
	"microService/modules/order/orderPb"
	orderUsecase "microService/modules/order/orderUsecase"
)

type (
	orderGrpcHandler struct {
		orderPb.UnimplementedOrderQueryServer
		orderUsecase orderUsecase.OrderUsecaseInterface
	}
)

func NewOrderGrpcHandler(orderUsecase orderUsecase.OrderUsecaseInterface) *orderGrpcHandler {
	return &orderGrpcHandler{orderUsecase: orderUsecase}
}
func (g *orderGrpcHandler) OrderQuery(ctx context.Context, req *orderPb.GetOrderRequest) (*orderPb.GetOrderResponse, error) {
	return nil, nil
}
