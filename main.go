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
	db := database.DbConn(ctx, &cfg)
	defer db.Disconnect(ctx)

	server.Start(ctx, &cfg, db)
}
