package paymentUsecase

import (
	"context"
	"microService/modules/payment/domain"
	"microService/modules/payment/paymentRepository"
	"microService/pkg/database"
	"time"

	"github.com/google/uuid"
)

type (
	PaymentUsecase interface {
		CreatePendingPayment(ctx context.Context, orderID, userID string, amount float64, currency string) (*domain.Payment, error)
	}

	paymentUsecase struct {
		paymentRepository paymentRepository.PaymentRepository
		txHelper          *database.TxHelper
	}
)

func NewPaymentUsecase(paymentRepository paymentRepository.PaymentRepository, tx *database.TxHelper) PaymentUsecase {
	return &paymentUsecase{
		paymentRepository: paymentRepository,
		txHelper:          tx,
	}
}
func (uc *paymentUsecase) CreatePendingPayment(ctx context.Context, orderID, userID string, amount float64, currency string) (*domain.Payment, error) {
	now := time.Now().UTC()
	p := &domain.Payment{
		ID:        "PAY-" + uuid.NewString(),
		OrderID:   orderID,
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		Status:    domain.PaymentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	trace := map[string]string{
		"source":    "payment-service",
		"component": "usecase",
		"trace_id":  uuid.NewString(),
	}

	if err := uc.txHelper.Transaction(ctx, func(ctx context.Context) error {
		return uc.paymentRepository.CreatePending(ctx, p, trace)
	}); err != nil {
		return nil, err
	}

	return p, nil
}
