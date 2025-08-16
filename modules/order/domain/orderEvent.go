package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID    string            `json:"event_id"`
	EventType  string            `json:"event_type"`
	Version    int               `json:"version"`
	Key        string            `json:"key"`
	OccurredAt time.Time         `json:"occurred_at"`
	Payload    any               `json:"payload"`
	Trace      map[string]string `json:"trace,omitempty"`
}

const (
	EventOrderCreated      = "order.created"
	EventOrderConfirmed    = "order.confirmed"
	EventOrderCancelled    = "order.cancelled"
	OutboxStatusPending    = "pending"
	OutboxStatusDispatched = "dispatched"
)

type Outbox struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Aggregate string    `gorm:"type:varchar(32);index;not null" json:"aggregate"` // "order"
	EventType string    `gorm:"type:varchar(64);index;not null" json:"event_type"`
	Key       string    `gorm:"index;not null" json:"key"` // partition/orderId
	Payload   []byte    `gorm:"type:jsonb;not null" json:"payload"`
	Status    string    `gorm:"type:varchar(16);index;not null" json:"status"` // pending, dispatched
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Outbox) TableName() string { return "outbox" }

func NewOrderEvent(eventType, key string, payload any, trace map[string]string) Event {
	return Event{
		EventID:    uuid.NewString(),
		EventType:  eventType,
		Version:    1,
		Key:        key,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
		Trace:      trace,
	}
}

func (e Event) ToOutbox(aggregate string) (Outbox, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return Outbox{}, err
	}
	return Outbox{
		Aggregate: aggregate,
		EventType: e.EventType,
		Key:       e.Key,
		Payload:   b,
		Status:    OutboxStatusPending,
	}, nil
}
