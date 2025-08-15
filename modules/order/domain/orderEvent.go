package domain

import "time"

type Event struct {
	EventID    string
	EventType  string
	Version    int
	Key        string // partition key (orderId)
	OccurredAt time.Time
	Payload    any
	Trace      map[string]string
}

// Event types for Order
const (
	EventOrderCreated   = "order.created"
	EventOrderConfirmed = "order.confirmed"
	EventOrderCancelled = "order.cancelled"
)

// Outbox record
type Outbox struct {
	ID        int64
	Aggregate string // "order"
	EventType string
	Key       string // orderId
	Payload   []byte
	Status    string // pending, dispatched
	CreatedAt time.Time
	UpdatedAt time.Time
}
