package itemRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepository interface{}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewitemRepository(db *mongo.Client) ItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item")
}
