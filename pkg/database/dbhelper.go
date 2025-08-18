package database

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func GetDB(ctx context.Context, base *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return base
}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

type TxHelper struct {
	DB DatabasesPostgres 
}

func NewTxHelper(db DatabasesPostgres) *TxHelper {
	return &TxHelper{DB: db}
}

func (h *TxHelper) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	base := h.DB.Connect()
	return base.Transaction(func(tx *gorm.DB) error {
		ctx = WithTx(ctx, tx)
		return fn(ctx)
	})
}
