package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID        string      `gorm:"primaryKey;type:uuid;not null" json:"id"`
	UserID    string      `gorm:"index;not null" json:"user_id"`
	Amount    int64       `gorm:"not null" json:"amount"`                // minor unit (satang)
	Currency  string      `gorm:"type:char(3);not null" json:"currency"` // THB
	Status    OrderStatus `gorm:"type:varchar(16);not null" json:"status"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items"` // ใช้ relation
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID string `gorm:"index;not null" json:"order_id"`
	SKU     string `gorm:"index;not null" json:"sku"`
	Qty     int32  `gorm:"not null" json:"qty"`
	Price   int64  `gorm:"not null" json:"price"` // price per unit
	// CreatedAt/UpdatedAt ถ้าอยากเก็บก็เติมได้
}

func (OrderItem) TableName() string { return "order_items" }

type CreateOrderInput struct {
	UserID   string      `json:"user_id" validate:"required"`
	Items    []OrderItem `json:"items" validate:"min=1,dive"`
	Amount   int64       `json:"amount" validate:"gt=0"`
	Currency string      `json:"currency" validate:"required"`
	// Meta     map[string]any `json:"meta,omitempty"`
}
