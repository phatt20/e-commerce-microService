package sagaUsecase

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func hdrGet(hs []*sarama.RecordHeader, k string) string {
	for _, h := range hs {
		if string(h.Key) == k {
			return string(h.Value)
		}
	}
	return ""
}

func hdrMap(hs []*sarama.RecordHeader) map[string]string {
	m := make(map[string]string, len(hs))
	for _, h := range hs {
		m[string(h.Key)] = string(h.Value)
	}
	return m
}

func newCmdHeaders(parent map[string]string, ceType, orderID, sagaID string) map[string]string {
	h := map[string]string{
		"ce-type":        ceType,
		"correlation-id": orderID,
		"saga-id":        sagaID,
		"command-id":     uuid.NewString(),
		"content-type":   "application/json",
		"sent-at":        time.Now().UTC().Format(time.RFC3339Nano),
	}
	if v := parent["trace-id"]; v != "" {
		h["trace-id"] = v
	}
	if v := parent["span-id"]; v != "" {
		h["parent-id"] = v
	}
	if v := parent["ttl"]; v != "" {
		h["ttl"] = v
	}
	if v := parent["event-version"]; v != "" {
		h["event-version"] = v
	}
	return h
}
