package inventoryUsecase

import "microService/modules/inventory/inventoryRepository"

type (
	InventoryUsecase interface{}
	inventoryUsecase struct {
		inventoryRepository inventoryRepository.InventoryRepository
	}
)

func NewInventoryUsecase(inventoryRepository inventoryRepository.InventoryRepository) InventoryUsecase {
	return &inventoryUsecase{inventoryRepository}
}
