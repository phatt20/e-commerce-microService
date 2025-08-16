package database

import (
	"context"
	"log"
	"microService/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConn(pctx context.Context, cfg *config.Config) *mongo.Client {
	if cfg.Mongo == nil || cfg.Mongo.Url == "" {
		log.Println("⚠️ Mongo config not found, skipping Mongo connection...")
		return nil
	}

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Url))
	if err != nil {
		log.Fatalf("เชื่อมต่อ MongoDB error: %s", err.Error())
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping MongoDB error: %s", err.Error())
	}

	log.Println("✅ Connected to MongoDB")
	return client
}
