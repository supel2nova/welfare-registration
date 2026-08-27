package verifier

import "context"

type Stub struct {
	FailAll     bool
	Unavailable bool
}

func (s Stub) Verify(_ context.Context, _ Request) (Result, error) {
	if s.Unavailable {
		return Result{}, ErrUnavailable
	}
	if s.FailAll {
		return Result{Matched: false, ProviderCode: "STUB"}, nil
	}
	return Result{
		Matched:      true,
		ProviderCode: "STUB",
		ReferenceNo:  "STUB-REF",
	}, nil
}
