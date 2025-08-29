package server

import (
	"context"
	"log"
	"time"

	paymentHttpHandler "microService/modules/payment/paymentHandler/http"
	paymentQueueHandler "microService/modules/payment/paymentHandler/paymentQueue"
	"microService/modules/payment/paymentRepository"
	"microService/modules/payment/paymentUsecase"
	"microService/pkg/database"
	"microService/pkg/outbox"
)

func (s *server) paymentService() {
	ctx := context.Background()

	repo := paymentRepository.NewPaymentRepository(s.postgres)
	txHelper := database.NewTxHelper(s.postgres)
	outboxRepo := outbox.NewOutboxRepo(s.postgres)

	usecase := paymentUsecase.NewPaymentUsecase(repo, txHelper)

	httpHandler := paymentHttpHandler.NewPaymentHttpHandler(s.cfg, usecase)
	// grpcHandler := paymentGrpcHandler.NewpaymentGrpcHandler(usecase)
	queueHandler := paymentQueueHandler.NewpaymentQueueHandler(s.cfg, usecase)

	_ = queueHandler
	_ = httpHandler

	// go func() {
	// 	grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
	// 	paymentPb.RegisterPaymentServiceServer(grpcServer, grpcHandler)
	// 	log.Printf("Payment gRPC server listening on %s", s.cfg.Grpc.PaymentUrl)
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		log.Println("gRPC server stopped:", err)
	// 	}
	// }()

	// Start Outbox worker
	worker := outbox.NewOutboxPublisher(
		outboxRepo,
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
		"payment",
		3*time.Second,
	)
	go func() {
		if err := worker.Run(ctx); err != nil {
			log.Println("Outbox worker stopped:", err)
		}
	}()

	payment := s.app.Group("/payment_v1")
	payment.GET("/health", s.healthCheckService)

}
