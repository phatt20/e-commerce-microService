package sagaQueue

import (
	"log"
	"microService/modules/saga/sagaUsecase"

	"github.com/IBM/sarama"
)

type SagaQueueHandler struct {
	uc sagaUsecase.SagaUsecase
}

func NewSagaQueueHandler(uc sagaUsecase.SagaUsecase) *SagaQueueHandler {
	return &SagaQueueHandler{uc: uc}
}

func (h *SagaQueueHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *SagaQueueHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *SagaQueueHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		evtType := header(msg.Headers, "ce-type")
		orderID := header(msg.Headers, "correlation-id")
		sagaID := header(msg.Headers, "saga-id")

		if err := h.uc.Handle(evtType, orderID, sagaID, msg.Value, msg.Headers); err != nil {
			log.Printf("saga handle error: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func header(headers []*sarama.RecordHeader, key string) string {
	for _, h := range headers {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}
