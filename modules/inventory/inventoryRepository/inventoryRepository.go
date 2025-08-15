package inventoryRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	InventoryRepository interface{}
	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) InventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory")
}
