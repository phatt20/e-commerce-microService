package paymentRepository

import (
	"context"
	"encoding/json"
	"microService/modules/payment/domain"
	"microService/pkg/database"
	"microService/pkg/outbox"
)

type (
	PaymentRepository interface {
		CreatePending(ctx context.Context, p *domain.Payment, sagaID string, trace map[string]string) error
		UpdateStatus(ctx context.Context, paymentID, status string, trace map[string]string) error
	}

	paymentRepository struct {
		db database.DatabasesPostgres
	}
)

func NewPaymentRepository(db database.DatabasesPostgres) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePending(ctx context.Context, p *domain.Payment, sagaID string, trace map[string]string) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)

	if err := db.Create(p).Error; err != nil {
		return err
	}

	ev := outbox.NewEvent("payment.pending", p.OrderID, map[string]any{
		"payment_id": p.ID,
		"order_id":   p.OrderID,
		"user_id":    p.UserID,
		"amount":     p.Amount,
		"currency":   p.Currency,
		"status":     p.Status,
	}, map[string]string{
		"ce-type":        "payment.pending",
		"correlation-id": p.OrderID,
		"content-type":   "application/json",
		"saga-id":        sagaID,
	}, trace)

	ob, err := ev.ToOutboxForPaymentPending("payment")
	if err != nil {
		return err
	}

	return db.Table("outbox").Create(&ob).Error
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, paymentID, status string, trace map[string]string) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)

	if err := db.Model(&domain.Payment{}).Where("id = ?", paymentID).Update("status", status).Error; err != nil {
		return err
	}

	var eventType string
	switch status {
	case domain.PaymentStatusSuccess:
		eventType = "payment.success"
	case domain.PaymentStatusFailed:
		eventType = "payment.failed"
	default:
		return nil
	}

	// outbox เดิมของ payment นี้
	var ob outbox.Outbox
	if err := db.Table("outbox").
		Where("key = ? AND aggregate = ?", paymentID, "payment").
		Order("created_at DESC").
		First(&ob).Error; err != nil {
		return err
	}

	// 4. อัพเดท eventType และ status ของ outbox ให้ pending
	ob.EventType = eventType
	ob.Status = outbox.OutboxStatusPending
	ob.Payload, _ = json.Marshal(map[string]any{
		"payment_id": paymentID,
		"status":     status,
		"trace":      trace,
	})

	// 5. บันทึกกลับ
	return db.Table("outbox").Save(&ob).Error
}
