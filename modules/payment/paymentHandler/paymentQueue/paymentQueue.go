package paymentqueue

import (
	"microService/config"
	"microService/modules/payment/paymentUsecase"
)

type (
	PaymentQueueHandler interface{}

	paymentQueueHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.PaymentUsecase
	}
)

func NewpaymentQueueHandler(cfg *config.Config, paymentUsecase paymentUsecase.PaymentUsecase) PaymentQueueHandler {
	return &paymentQueueHandler{cfg, paymentUsecase}
}
//ยัังไม่เจอที่ต้องทำ เเละยังไม่รู้จะดักเอา messagse อะไรเข้ามา