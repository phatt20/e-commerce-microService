package itemHandler

import (
	"microService/config"
	"microService/modules/item/itemUsecase"
)

type (
	ItemHttpHandler interface{}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.ItemUsecase
	}
)

func NewitemHttpHandler(cfg *config.Config, itemUsecase itemUsecase.ItemUsecase) ItemHttpHandler {
	return &itemHttpHandler{cfg, itemUsecase}
}
