package orderRepo

import (
	"context"
	"log"
	"microService/modules/order/domain"
	"microService/modules/payment/paymentPb"
	"microService/pkg/database"
	grpccon "microService/pkg/grpcCon"
	"microService/pkg/jwtauth"
	"microService/pkg/outbox"
	"time"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, o *domain.Order) error
	CreateOrderWithOutbox(ctx context.Context, o *domain.Order, trace map[string]string) error
}

type Repo struct {
	db database.DatabasesPostgres // Connect() *gorm.DB
	tx *database.TxHelper
}

func NewRepo(db database.DatabasesPostgres) *Repo {
	return &Repo{
		db: db,
		tx: database.NewTxHelper(db),
	}
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

func (r *Repo) CreateOrderWithOutbox(ctx context.Context, paymentUrl string, order *domain.Order, trace map[string]string) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)

	// 1️⃣ สร้าง order
	if err := db.Create(order).Error; err != nil {
		return err
	}

	// 2️⃣ สร้าง event สำหรับ outbox
	ev := outbox.NewEvent(domain.EventOrderCreated, order.ID, map[string]any{
		"order_id": order.ID,
		"user_id":  order.UserID,
		"amount":   order.Amount,
		"currency": order.Currency,
		"status":   order.Status,
		"items":    order.Items,
	}, trace)

	outboxRow, err := ev.ToOutbox("order")
	if err != nil {
		return err
	}

	if err := db.Table("outbox").Create(&outboxRow).Error; err != nil {
		return err
	}

	if err := r.callPayment(ctx, paymentUrl, order); err != nil {
		log.Println("Payment gRPC failed, order still created:", err)
	}

	return nil
}

// grpc payment naja
func (r *Repo) callPayment(ctx context.Context, paymentUrl string, order *domain.Order) error {

	pctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&pctx)

	conn, err := grpccon.NewGrpcClient(paymentUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err)
		return err
	}

	amountFloat := float64(order.Amount) / 100
	req := &paymentPb.CreatePaymentRequest{
		OrderId:  order.ID,
		UserId:   order.UserID,
		Amount:   amountFloat,
		Currency: order.Currency,
	}

	resp, err := conn.Payment().CreatePayment(pctx, req)
	if err != nil {
		log.Printf("Error: CreatePayment failed: %s", err)
		return err
	}

	log.Printf("Payment created: ID=%s, Status=%s", resp.GetPaymentId(), resp.GetStatus())
	return nil
}
