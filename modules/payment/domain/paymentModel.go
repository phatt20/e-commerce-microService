package domain

import "time"

const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)

type Payment struct {
	ID        string `gorm:"primaryKey;type:varchar(64)"`
	OrderID   string `gorm:"index;type:varchar(64)"`
	UserID    string `gorm:"type:varchar(64)"`
	Amount    float64
	Currency  string `gorm:"type:varchar(8)"`
	Status    string `gorm:"type:varchar(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Payment) TableName() string { return "payments" }
