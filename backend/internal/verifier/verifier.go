package verifier

import (
	"context"
	"errors"
)

var ErrUnavailable = errors.New("verifier: provider ไม่พร้อมให้บริการ")

type Request struct {
	NationalID string
	LaserID    string
}

type Result struct {
	Matched      bool
	ProviderCode string
	ReferenceNo  string
}

type Verifier interface {
	Verify(ctx context.Context, req Request) (Result, error)
}
