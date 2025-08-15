package server

import (
	paymentGrpcHandler "microService/modules/payment/paymentHandler/grpc"
	paymentHttpHandler "microService/modules/payment/paymentHandler/http"
	paymentQueueHandler "microService/modules/payment/paymentHandler/paymentQueue"

	"microService/modules/payment/paymentRepository"
	"microService/modules/payment/paymentUsecase"
)

func (s *server) paymentService() {
	repo := paymentRepository.NewPaymentRepository(s.db)
	usecase := paymentUsecase.NewPaymentUsecase(repo)
	paymentHttpHandler := paymentHttpHandler.NewPaymentHttpHandler(s.cfg, usecase)
	paymentGrpcHandler := paymentGrpcHandler.NewpaymentGrpcHandler(usecase)
	paymentQueueHandler := paymentQueueHandler.NewpaymentQueueHandler(s.cfg, usecase)

	_ = paymentHttpHandler
	_ = paymentGrpcHandler
	_ = paymentQueueHandler

	payment := s.app.Group("/payment_v1")

	payment.GET("/health", s.healthCheckService)
}
