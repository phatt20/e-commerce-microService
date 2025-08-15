package inventoryHandler

import (
	"microService/config"
	"microService/modules/inventory/inventoryUsecase"
)

type (
	InventoryHttpHandler interface{}
	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecase
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecase) InventoryHttpHandler {
	return &inventoryHttpHandler{cfg, inventoryUsecase}
}
