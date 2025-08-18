package orderhanlder

import (
	"context"
	"microService/config"
	"microService/modules/order/domain"
	orderUsecase "microService/modules/order/orderUsecase"
	"microService/pkg/request"
	"microService/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OrderHttpHandler interface {
	CreateOrder(c echo.Context) error
}

type orderHttpHandler struct {
	cfg          *config.Config
	orderUsecase orderUsecase.OrderUsecaseInterface
}

func NewOrderHttpHandler(cfg *config.Config, orderUsecase orderUsecase.OrderUsecaseInterface) OrderHttpHandler {
	return &orderHttpHandler{cfg, orderUsecase}
}

func (h *orderHttpHandler) CreateOrder(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(domain.CreateOrderInput)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.orderUsecase.CreateOrder(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
