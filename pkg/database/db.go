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
	ctx, cencel := context.WithTimeout(pctx, 10*time.Second)
	defer cencel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.Url))

	if err != nil {
		log.Fatalf("เชื่อมต่อ db err %s", err.Error())
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("pinging db err %s", err.Error())
	}
	return client
}
