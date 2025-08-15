package inventoryHandler

import (
	"microService/config"
	"microService/modules/inventory/inventoryUsecase"
)

type (
	InventoryQueueHandler interface{}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecase
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecase) InventoryQueueHandler {
	return &inventoryQueueHandler{cfg, inventoryUsecase}
}
