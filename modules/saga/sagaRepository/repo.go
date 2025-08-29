package sagaRepository

import (
	"context"
	"sync"
)


type StateRepo interface {
	Get(ctx context.Context, orderID string) (step string, sagaID string, err error)
	Next(ctx context.Context, orderID, sagaID, next string) error
	Fail(ctx context.Context, orderID, sagaID string) error
}

type MemoryRepo struct {
	mu   sync.Mutex
	data map[string]struct{ Step, SagaID string }
}

func NewMemoryRepo() StateRepo {
	return &MemoryRepo{data: make(map[string]struct{ Step, SagaID string })}
}

func (r *MemoryRepo) Get(ctx context.Context, orderID string) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if x, ok := r.data[orderID]; ok {
		return x.Step, x.SagaID, nil
	}
	return "START", "", nil
}

func (r *MemoryRepo) Next(ctx context.Context, orderID, sagaID, next string) error {
	r.mu.Lock()
	r.data[orderID] = struct{ Step, SagaID string }{Step: next, SagaID: sagaID}
	r.mu.Unlock()
	return nil
}

func (r *MemoryRepo) Fail(ctx context.Context, orderID, sagaID string) error {
	r.mu.Lock()
	r.data[orderID] = struct{ Step, SagaID string }{Step: "FAILED", SagaID: sagaID}
	r.mu.Unlock()
	return nil
}
