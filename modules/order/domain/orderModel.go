package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID        string
	UserID    string
	Amount    int64  // minor unit (e.g., satang)
	Currency  string // THB
	Status    OrderStatus
	Items     []OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	SKU   string
	Qty   int32
	Price int64 // price per unit, optional
}

type CreateOrderInput struct {
	UserID   string         `json:"-"`
	Items    []OrderItem    `json:"items" validate:"min=1,dive"`
	Amount   int64          `json:"amount" validate:"gt=0"`
	Currency string         `json:"currency" validate:"required"`
	Meta     map[string]any `json:"meta,omitempty"`
}

