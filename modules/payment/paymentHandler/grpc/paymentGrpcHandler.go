package paymentHandler

import "microService/modules/payment/paymentUsecase"

type (
	paymentGrpcHandler struct {
		paymentUsecase paymentUsecase.PaymentUsecase
	}
)

func NewpaymentGrpcHandler(paymentUsecase paymentUsecase.PaymentUsecase) *paymentGrpcHandler {
	return &paymentGrpcHandler{paymentUsecase}
}
