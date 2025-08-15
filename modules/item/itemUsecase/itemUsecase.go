package itemUsecase

import "microService/modules/item/itemRepository"

type (
	ItemUsecase interface{}

	itemUsecase struct {
		itemRepositoty itemRepository.ItemRepository
	}
)

func NewitemUsecase(itemRepositoty itemRepository.ItemRepository) ItemUsecase {
	return &itemUsecase{itemRepositoty}
}
