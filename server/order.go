package server

import (
	"context"
	"log"
	orderhanlder "microService/modules/order/orderHanlder"
	"microService/modules/order/orderRepo"
	orderUsecase "microService/modules/order/orderUsecase"
	"microService/pkg/database"
	"microService/pkg/outbox"
	"time"
)

func (s *server) orderService() {
	ctx := context.Background()

	// 2. Repository
	repo := orderRepo.NewRepo(s.postgres)
	txHelper := database.NewTxHelper(s.postgres)

	outboxRepo := outbox.NewOutboxRepo(s.postgres)

	usecase := orderUsecase.NewOrderUsecase(repo, txHelper)
	orderHandler := orderhanlder.NewOrderHttpHandler(s.cfg, usecase)
	// grpcHandler := orderGrpcHandler.NewOrderGrpcHandler(usecase)
	// queueHandler := orderQueueHandler.NewOrderQueueHandler(s.cfg, usecase)

	// 5. Outbox worker (background)
	worker := outbox.NewOutboxPublisher(
		outboxRepo,
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
		"order",
		3*time.Second, // poll interval
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
