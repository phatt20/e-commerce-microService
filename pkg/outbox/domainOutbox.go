package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	OutboxStatusPending    = "pending"
	OutboxStatusDispatched = "dispatched"
)

type Event struct {
	EventID    string            `json:"event_id"`
	EventType  string            `json:"event_type"`
	Version    int               `json:"version"`
	Key        string            `json:"key"`
	OccurredAt time.Time         `json:"occurred_at"`
	Payload    any               `json:"payload"`
	Trace      map[string]string `json:"trace,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type Outbox struct {
	ID         int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Aggregate  string            `gorm:"type:varchar(32);index;not null" json:"aggregate"`
	EventID    string            `gorm:"type:varchar(36);index;not null" json:"event_id"`
	EventType  string            `gorm:"type:varchar(64);index;not null" json:"event_type"`
	Key        string            `gorm:"index;not null" json:"key"` // partition/orderId
	OccurredAt time.Time         `gorm:"not null" json:"occurred_at"`
	Payload    []byte            `gorm:"type:jsonb;not null" json:"payload"`
	Headers    map[string]string `gorm:"type:jsonb;not null" json:"headers"`
	Trace      map[string]string `gorm:"type:jsonb" json:"trace,omitempty"`
	Status     string            `gorm:"type:varchar(16);index;not null" json:"status"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Outbox) TableName() string { return "outbox" }

func NewEvent(eventType, key string, payload any, trace, headers map[string]string) Event {
	return Event{
		EventID:    uuid.NewString(),
		EventType:  eventType,
		Version:    1,
		Key:        key,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
		Trace:      trace,
		Headers:    headers,
	}
}

func (e Event) ToOutbox(aggregate string) (Outbox, error) {
	b, err := json.Marshal(e.Payload)
	if err != nil {
		return Outbox{}, err
	}

	return Outbox{
		Aggregate:  aggregate,
		EventID:    e.EventID,
		EventType:  e.EventType,
		Key:        e.Key,
		OccurredAt: e.OccurredAt,
		Payload:    b,
		Headers:    e.Headers,
		Trace:      e.Trace,
		Status:     OutboxStatusPending,
	}, nil
}

func (e Event) ToOutboxForPaymentPending(aggregate string) (Outbox, error) {
	b, err := json.Marshal(e.Payload)
	if err != nil {
		return Outbox{}, err
	}

	return Outbox{
		Aggregate:  aggregate,
		EventID:    e.EventID,
		EventType:  e.EventType,
		Key:        e.Key,
		OccurredAt: e.OccurredAt,
		Payload:    b,
		Headers:    e.Headers,
		Trace:      e.Trace,
		Status:     OutboxStatusDispatched,
	}, nil
}
