package paymentHandler

import (
	"microService/config"
	"microService/modules/payment/paymentUsecase"
)

type (
	PaymentHttpHandler interface{}

	paymentHttpHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.PaymentUsecase
	}
)

func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase paymentUsecase.PaymentUsecase) PaymentHttpHandler {
	return &paymentHttpHandler{cfg, paymentUsecase}
}
//ยัังไม่เจอที่ต้องทำ