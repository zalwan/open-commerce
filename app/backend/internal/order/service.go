// Package order owns checkout and order lifecycle.
package order

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/payment"
)

// Status lifecycle: pending -> paid -> shipped -> done; payment_failed/cancelled terminal from early states.
type Status string

const (
	StatusPending       Status = "pending"
	StatusPaid          Status = "paid"
	StatusPaymentFailed Status = "payment_failed"
	StatusShipped       Status = "shipped"
	StatusDone          Status = "done"
	StatusCancelled     Status = "cancelled"
)

// Order is a checkout snapshot.
type Order struct {
	ID         string      `json:"id"`
	SessionID  string      `json:"sessionId"`
	Email      string      `json:"email"`
	Items      []cart.Item `json:"items"`
	TotalMinor int64       `json:"totalMinor"`
	Currency   string      `json:"currency"`
	Status     Status      `json:"status"`
	PaymentTx  string      `json:"paymentTx,omitempty"`
	CreatedAt  time.Time   `json:"createdAt"`
}

var (
	ErrEmptyCart  = errors.New("cart is empty")
	ErrBadEmail   = errors.New("email is required")
	ErrNotFound   = errors.New("order not found")
	ErrBadStatus  = errors.New("invalid status transition")
	ErrChargeFail = errors.New("payment failed")
)

// Service coordinates cart snapshot + payment stub + order store.
type Service struct {
	mu      sync.Mutex
	orders  map[string]*Order
	seq     atomic.Int64
	carts   *cart.Service
	payment payment.Provider
}

func NewService(carts *cart.Service, pay payment.Provider) *Service {
	return &Service{orders: make(map[string]*Order), carts: carts, payment: pay}
}

// Checkout snapshots the cart, charges via stub, records order, clears cart on success.
func (s *Service) Checkout(sessionID, email string, failPayment bool) (Order, error) {
	if email == "" {
		return Order{}, ErrBadEmail
	}
	c := s.carts.Get(sessionID)
	if len(c.Items) == 0 {
		return Order{}, ErrEmptyCart
	}
	txID, err := s.payment.Charge(c.Subtotal, "IDR", failPayment)
	status := StatusPaid
	if err != nil {
		status = StatusPaymentFailed
	}
	n := s.seq.Add(1)
	o := Order{
		ID:         fmt.Sprintf("o-%d", n),
		SessionID:  sessionID,
		Email:      email,
		Items:      c.Items,
		TotalMinor: c.Subtotal,
		Currency:   "IDR",
		Status:     status,
		PaymentTx:  txID,
		CreatedAt:  time.Now().UTC(),
	}
	s.mu.Lock()
	s.orders[o.ID] = &o
	s.mu.Unlock()
	if err != nil {
		return o, ErrChargeFail
	}
	s.carts.Clear(sessionID)
	return o, nil
}

func (s *Service) Get(id string) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	return *o, nil
}

func (s *Service) List() []Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, *o)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}

// SetStatus enforces a linear lifecycle.
func (s *Service) SetStatus(id string, next Status) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	allowed := map[Status][]Status{
		StatusPending:       {StatusPaid, StatusPaymentFailed, StatusCancelled},
		StatusPaid:          {StatusShipped, StatusCancelled},
		StatusPaymentFailed: {StatusCancelled},
		StatusShipped:       {StatusDone},
		StatusDone:          {},
		StatusCancelled:     {},
	}
	// Paid-from-stub orders start at paid; allow checkout-created paid to ship.
	if o.Status == StatusPaid {
		allowed[StatusPaid] = []Status{StatusShipped, StatusCancelled}
	}
	okTransition := false
	for _, n := range allowed[o.Status] {
		if n == next {
			okTransition = true
			break
		}
	}
	if !okTransition {
		return *o, ErrBadStatus
	}
	o.Status = next
	return *o, nil
}
