package paymentUsecase

import "microService/modules/payment/paymentRepository"

type (
	PaymentUsecase interface{}

	paymentUsecase struct {
		paymentRepository paymentRepository.PaymentRepository
	}
)

func NewPaymentUsecase(paymentRepository paymentRepository.PaymentRepository) PaymentUsecase {
	return &paymentUsecase{paymentRepository}
}
