// Package payment is a deterministic stub. Never processes real money.
package payment

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// Provider abstracts charge/refund so a real gateway can replace the stub.
type Provider interface {
	Charge(amountMinor int64, currency string, fail bool) (txID string, err error)
	Refund(txID string) error
}

var ErrChargeFailed = errors.New("charge failed (stub)")

// Stub always succeeds unless amount invalid or fail=true (for testing failure path).
type Stub struct {
	counter atomic.Int64
}

func NewStub() *Stub { return &Stub{} }

func (s *Stub) Charge(amountMinor int64, currency string, fail bool) (string, error) {
	if amountMinor <= 0 {
		return "", ErrChargeFailed
	}
	if fail {
		return "", ErrChargeFailed
	}
	n := s.counter.Add(1)
	return fmt.Sprintf("stub-tx-%d", n), nil
}

func (s *Stub) Refund(txID string) error {
	if txID == "" {
		return errors.New("txID required")
	}
	return nil
}
