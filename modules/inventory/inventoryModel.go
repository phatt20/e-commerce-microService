package inventory

import (
	"microService/modules/item"
	"microService/modules/models"
)

type (
	UpdateInventoryReq struct {
		UserId string `json:"user_id" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInInventory struct {
		InventoryId string `json:"inventory_id"`
		UserId    string `json:"user_id"`
		*item.ItemShowCase
	}

	InventorySearchReq struct {
		models.PaginateReq
	}

	RollbackUserInventoryReq struct {
		InventoryId string `json:"inventory_id"`
		UserId    string `json:"user_id"`
		ItemId      string `json:"item_id"`
	}
)
