package sagaUsecase

import (
	"context"
	"encoding/json"
	"microService/modules/saga/sagaRepository"
	"microService/pkg/queue"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type SagaUsecase interface {
	Handle(evtType, orderID, sagaID string, payload []byte, headers []*sarama.RecordHeader) error
}
type Topics struct {
	OrderCmd     string
	PaymentCmd   string
	InventoryCmd string
	ShippingCmd  string
}

type usecase struct {
	state   sagaRepository.StateRepo
	idem    sagaRepository.IdemRepo
	topics  Topics
	brokers []string
	apiKey  string
	secret  string
}

func New(state sagaRepository.StateRepo, idem sagaRepository.IdemRepo, topics Topics, brokers []string, apiKey, secret string) SagaUsecase {
	return &usecase{state: state, idem: idem, topics: topics, brokers: brokers, apiKey: apiKey, secret: secret}
}

func (u *usecase) Handle(evtType, orderID, sagaID string, _ []byte, hdrs []*sarama.RecordHeader) error {
	parent := hdrMap(hdrs)
	if sagaID == "" {
		_, cur, _ := u.state.Get(context.Background(), orderID)
		if cur != "" {
			sagaID = cur
		} else {
			sagaID = "SAGA-" + uuid.NewString()
		}
	}
	if id := parent["event-id"]; id != "" {
		done, _ := u.idem.WasProcessed(context.Background(), id)
		if done {
			return nil
		}
		defer u.idem.MarkProcessed(context.Background(), id)
	}

	switch evtType {
	case "order.created":
		u.state.Next(context.Background(), orderID, sagaID, "RESERVE_STOCK")
		u.cmd(u.topics.InventoryCmd, orderID, "inventory.reserve", sagaID, parent, nil)

	case "inventory.reserved":
		u.state.Next(context.Background(), orderID, sagaID, "AUTHORIZE_PAYMENT")
		u.cmd(u.topics.PaymentCmd, orderID, "payment.authorize", sagaID, parent, nil)

	case "payment.authorized":
		u.state.Next(context.Background(), orderID, sagaID, "CONFIRM_RESERVATION")
		u.cmd(u.topics.InventoryCmd, orderID, "inventory.confirm", sagaID, parent, nil)

	case "inventory.deducted":
		u.state.Next(context.Background(), orderID, sagaID, "SCHEDULE_SHIPMENT")
		u.cmd(u.topics.ShippingCmd, orderID, "shipping.schedule", sagaID, parent, nil)

	case "shipment.scheduled":
		u.state.Next(context.Background(), orderID, sagaID, "CONFIRM_ORDER")
		u.cmd(u.topics.OrderCmd, orderID, "order.confirm", sagaID, parent, nil)

	case "inventory.released":
		u.state.Fail(context.Background(), orderID, sagaID)
		u.cmd(u.topics.OrderCmd, orderID, "order.cancel", sagaID, parent, map[string]any{"reason": "reserve_failed"})

	case "payment.failed":
		u.state.Fail(context.Background(), orderID, sagaID)
		u.cmd(u.topics.InventoryCmd, orderID, "inventory.release", sagaID, parent, nil)
		u.cmd(u.topics.OrderCmd, orderID, "order.cancel", sagaID, parent, map[string]any{"reason": "payment_failed"})

	case "shipment.failed":
		u.state.Fail(context.Background(), orderID, sagaID)
		u.cmd(u.topics.PaymentCmd, orderID, "payment.refund", sagaID, parent, nil)
		u.cmd(u.topics.OrderCmd, orderID, "order.cancel", sagaID, parent, map[string]any{"reason": "shipment_failed"})
	}
	return nil
}

func (u *usecase) cmd(topic, orderID, ceType, sagaID string, parent map[string]string, extra map[string]any) error {
	body := map[string]any{"order_id": orderID}
	for k, v := range extra {
		body[k] = v
	}
	b, _ := json.Marshal(body)
	h := newCmdHeaders(parent, ceType, orderID, sagaID)
	return queue.PushMessageWithKeyAndHeadersToQueue(u.brokers, u.apiKey, u.secret, topic, orderID, b, h)
}
