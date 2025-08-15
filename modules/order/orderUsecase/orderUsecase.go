package usecase

import (
	"context"
	"encoding/json"
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
	repo      orderRepo.OrderRepository
	outbox    orderRepo.OutboxRepository
	validator *validator.Validate
}

func NewOrderUsecase(r orderRepo.OrderRepository, ob orderRepo.OutboxRepository) OrderUsecaseInterface {
	return &OrderUsecase{repo: r, outbox: ob, validator: validator.New()}
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

	// within ONE DB transaction: write order + outbox
	err := uc.repo.Tx(ctx, func(ctx context.Context) error {
		if err := uc.repo.InsertOrder(ctx, order); err != nil {
			return err
		}
		payload := map[string]any{
			"orderId":   order.ID,
			"userId":    order.UserID,
			"items":     order.Items,
			"amount":    order.Amount,
			"currency":  order.Currency,
			"createdAt": now,
		}
		b, _ := json.Marshal(payload)
		return uc.outbox.Add(ctx, &domain.Outbox{
			Aggregate: "order",
			EventType: domain.EventOrderCreated,
			Key:       order.ID,
			Payload:   b,
			Status:    "pending",
			CreatedAt: now,
			UpdatedAt: now,
		})
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
