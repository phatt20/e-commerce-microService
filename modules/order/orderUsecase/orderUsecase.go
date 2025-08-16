package usecase

import (
	"context"
	"microService/modules/order/domain"
	"microService/modules/order/orderRepo"
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
}

func NewOrderUsecase(r *orderRepo.Repo) OrderUsecaseInterface {
	return &OrderUsecase{
		repo:      r,
		validator: validator.New(),
	}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, in *domain.CreateOrderInput) (*domain.Order, error) {
	if err := uc.validator.Struct(in); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	orderID := "ORD-" + uuid.NewString()

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
		"source": "usecase",
	}

	if err := uc.repo.Tx(ctx, func(ctx context.Context) error {
		return uc.repo.CreateOrderWithOutbox(ctx, order, trace)
	}); err != nil {
		return nil, err
	}

	return order, nil
}
