package usecase

import (
	"context"
	"microService/modules/order/domain"
	"microService/modules/order/orderRepo"
	"microService/pkg/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OrderUsecaseInterface interface {
	CreateOrder(ctx context.Context, in *domain.CreateOrderInput) (*domain.Order, error)
}

type OrderUsecase struct {
	repo      *orderRepo.Repo
	validator *validator.Validate
	txHelper  *database.TxHelper
}

func NewOrderUsecase(r *orderRepo.Repo, tx *database.TxHelper) OrderUsecaseInterface {
	return &OrderUsecase{
		repo:      r,
		txHelper:  tx,
		validator: validator.New(),
	}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, in *domain.CreateOrderInput) (*domain.Order, error) {
	if err := uc.validator.Struct(in); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	orderID := uuid.NewString()

	order := &domain.Order{
		ID:        orderID,
		UserID:    in.UserID,
		Amount:    in.Amount,
		Currency:  in.Currency,
		Status:    domain.OrderStatusPending,
		Items:     in.Items,
		CreatedAt: now,
		UpdatedAt: now,
	}

	trace := map[string]string{
		"source":    "order-service",
		"component": "usecase",
		"trace_id":  uuid.NewString(),
	}

	if err := uc.txHelper.Transaction(ctx, func(ctx context.Context) error {
		return uc.repo.CreateOrderWithOutbox(ctx, order, trace)
	}); err != nil {
		return nil, err
	}

	return order, nil
}
