package payment

type (
	ItemServiceReq struct {
		Items []*ItemServiceReqDatum `json:"items" validate:"required"`
	}

	ItemServiceReqDatum struct {
		ItemId string  `json:"item_id" validate:"required,max=64"`
		Price  float64 `json:"price" validate:"required,gt=0"`
	}

	PaymentTransferReq struct {
		UserId string  `json:"user_id" validate:"required"`
		ItemId string  `json:"item_id" validate:"required,max=64"`
		Amount float64 `json:"amount" validate:"required,gt=0"`
	}

	PaymentTransferRes struct {
		InventoryId   string  `json:"inventory_id"`
		TransactionId string  `json:"transaction_id"`
		UserId        string  `json:"user_id"`
		ItemId        string  `json:"item_id"`
		Amount        float64 `json:"amount"`
		Error         string  `json:"error,omitempty"`
	}
)
