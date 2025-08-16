package main

import (
	"context"
	"log"
	"microService/config"
	"microService/pkg/database"
	"microService/server"
	"os"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	// Mongo
	dbM := database.DbConn(ctx, &cfg)
	if dbM != nil {
		defer func() {
			if err := dbM.Disconnect(ctx); err != nil {
				log.Printf("⚠️ error disconnecting Mongo: %v", err)
			}
		}()
	}

	// Postgres
	var dbPost database.DatabasesPostgres
	if cfg.Postgres != nil {
		dbPost = database.NewPostgresDatabase(cfg.Postgres)
	}

	// Run server
	server.Start(ctx, &cfg, dbPost, dbM)
}
