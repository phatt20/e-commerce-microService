package outbox

import (
	"context"
	"encoding/json"
	"log"
	"microService/pkg/queue"
	"time"
)

type Outbox struct {
	ID        int64
	EventType string
	Key       string
	Payload   string
}

type OutboxPublisher struct {
	outboxRepo OutboxRepository
	brokerURLs []string
	apiKey     string
	secret     string
	topic      string
	poll       time.Duration
}

func NewOutboxPublisher(
	r OutboxRepository,
	brokerURLs []string,
	apiKey string,
	secret string,
	topic string,
	every time.Duration,
) *OutboxPublisher {
	return &OutboxPublisher{
		outboxRepo: r,
		brokerURLs: brokerURLs,
		apiKey:     apiKey,
		secret:     secret,
		topic:      topic,
		poll:       every,
	}
}

func (w *OutboxPublisher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			rows, err := w.outboxRepo.PollPending(ctx, 100)
			if err != nil {
				log.Println("outbox poll error:", err)
				continue
			}

			for _, ob := range rows {
				env := map[string]any{
					"eventId":    ob.ID,
					"eventType":  ob.EventType,
					"version":    1,
					"occurredAt": time.Now().UTC().Format(time.RFC3339Nano),
					"key":        ob.Key,
					"payload":    json.RawMessage(ob.Payload),
				}

				b, _ := json.Marshal(env)

				err := queue.PushMessageWithKeyToQueue(
					w.brokerURLs,
					w.apiKey,
					w.secret,
					w.topic,
					ob.Key,
					b,
				)
				if err != nil {
					log.Println("publish error:", err)
					continue
				}

				if err := w.outboxRepo.MarkDispatched(ctx, ob.ID); err != nil {
					log.Println("mark dispatched error:", err)
				}
			}
		}
	}
}
