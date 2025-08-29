package sagaRepository

import (
	"context"
	"sync"
)

type IdemRepo interface {
	WasProcessed(ctx context.Context, id string) (bool, error)
	MarkProcessed(ctx context.Context, id string) error
}

type MemoryIdemRepo struct {
	mu sync.Mutex
	s  map[string]struct{}
}

func NewMemoryIdemRepo() IdemRepo { return &MemoryIdemRepo{s: map[string]struct{}{}} }

func (r *MemoryIdemRepo) WasProcessed(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.s[id]
	return ok, nil
}

func (r *MemoryIdemRepo) MarkProcessed(ctx context.Context, id string) error {
	r.mu.Lock()
	r.s[id] = struct{}{}
	r.mu.Unlock()
	return nil
}
