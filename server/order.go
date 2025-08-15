package server

import (
	"context"
	"log"
	orderhanlder "microService/modules/order/orderHanlder"
	"microService/modules/order/orderRepo"
	orderUsecase "microService/modules/order/orderUsecase"

	outboxP "microService/modules/order/outboxPush"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *server) orderService() {
	ctx := context.Background()

	// 1. DB pool
	pool, err := pgxpool.New(ctx, s.cfg.Db.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// 2. Repos
	repo := orderRepo.NewRepo(pool)
	outbox := orderRepo.NewOutboxRepo(pool)

	// 3. Usecase
	usecase := orderUsecase.NewOrderUsecase(repo, outbox)

	// 4. Handlers
	orderHandler := orderhanlder.NewOrderHttpHandler(s.cfg, usecase)
	// grpcHandler := orderGrpcHandler.NewOrderGrpcHandler(usecase)
	// queueHandler := orderQueueHandler.NewOrderQueueHandler(s.cfg, usecase)

	// 5. Outbox worker
	worker := outboxP.NewOutboxPublisher(
		outbox,
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
		"order", 3,
	)
	go func() {
		if err := worker.Run(ctx); err != nil {
			log.Println("outbox worker stopped:", err)
		}
	}()

	// 6. HTTP routes
	order := s.app.Group("/order_v1")
	order.POST("/order/create", orderHandler.CreateOrder, s.middleware.JwtAuthorization)
}
