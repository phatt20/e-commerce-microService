package paymentRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PaymentRepository interface{}
	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment")
}
