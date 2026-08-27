package verifier

import (
	"context"
	"sync"
)

type Spy struct {
	Inner Verifier
	mu    sync.Mutex
	Calls int
}

func (s *Spy) Verify(ctx context.Context, req Request) (Result, error) {
	s.mu.Lock()
	s.Calls++
	s.mu.Unlock()
	return s.Inner.Verify(ctx, req)
}
