package orderRepo

import (
	"context"
	"microService/modules/order/domain"
	"microService/pkg/database"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) error
	InsertOrder(ctx context.Context, o *domain.Order) error
	CreateOrderWithOutbox(ctx context.Context, o *domain.Order, trace map[string]string) error
}

type Repo struct {
	db database.DatabasesPostgres // Connect() *gorm.DB
}

func NewRepo(db database.DatabasesPostgres) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	base := r.db.Connect()
	return base.Transaction(func(tx *gorm.DB) error {
		ctx = database.WithTx(ctx, tx)
		return fn(ctx)
	})
}

func (r *Repo) InsertOrder(ctx context.Context, o *domain.Order) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)

	if err := db.Create(o).Error; err != nil {
		return err
	}
	if len(o.Items) > 0 {
		for i := range o.Items {
			o.Items[i].OrderID = o.ID
		}
		if err := db.CreateInBatches(o.Items, 100).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) CreateOrderWithOutbox(ctx context.Context, order *domain.Order, trace map[string]string) error {
	
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)

	
	if err := db.Create(order).Error; err != nil {
		return err 
	}

	
	ev := domain.NewOrderEvent(domain.EventOrderCreated, order.ID, map[string]any{
		"order_id": order.ID,
		"user_id":  order.UserID,
		"amount":   order.Amount,
		"currency": order.Currency,
		"status":   order.Status,
		"items":    order.Items,
	}, trace)

	outbox, err := ev.ToOutbox("order")
	if err != nil {
		return err
	}

	if err := db.Table("outbox").Create(&outbox).Error; err != nil {
		return err // fail → rollback
	}

	return nil
}
