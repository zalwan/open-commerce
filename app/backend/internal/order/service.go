// Package order owns checkout and order lifecycle.
package order

import (
	"context"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
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

// Store abstracts order persistence: memory (dev) or Postgres (DATABASE_URL set).
type Store interface {
	NextID(ctx context.Context) (string, error)
	Save(ctx context.Context, o Order) error
	Get(ctx context.Context, id string) (Order, error)
	List(ctx context.Context) ([]Order, error)
}

// MemoryStore keeps orders process-local with an atomic sequence.
type MemoryStore struct {
	mu     sync.Mutex
	orders map[string]*Order
	seq    atomic.Int64
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{orders: make(map[string]*Order)} }

func (s *MemoryStore) NextID(_ context.Context) (string, error) {
	return "o-" + itoa(s.seq.Add(1)), nil
}

func (s *MemoryStore) Save(_ context.Context, o Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := o
	s.orders[o.ID] = &cp
	return nil
}

func (s *MemoryStore) Get(_ context.Context, id string) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	return *o, nil
}

func (s *MemoryStore) List(_ context.Context) ([]Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, *o)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

// Service coordinates cart snapshot + stock reservation + payment stub + order store.
type Service struct {
	store   Store
	carts   *cart.Service
	catalog *catalog.Service
	payment payment.Provider
}

func NewService(store Store, carts *cart.Service, catalogSvc *catalog.Service, pay payment.Provider) *Service {
	return &Service{store: store, carts: carts, catalog: catalogSvc, payment: pay}
}

// Checkout reserves stock first, then charges: success records a paid order,
// payment failure compensates (restocks) and records payment_failed.
func (s *Service) Checkout(ctx context.Context, sessionID, email string, failPayment bool) (Order, error) {
	if email == "" {
		return Order{}, ErrBadEmail
	}
	c, err := s.carts.Get(ctx, sessionID)
	if err != nil {
		return Order{}, err
	}
	if len(c.Items) == 0 {
		return Order{}, ErrEmptyCart
	}
	// Reserve stock before charging so concurrent checkouts cannot oversell.
	reserved := make([]cart.Item, 0, len(c.Items))
	for _, it := range c.Items {
		if derr := s.catalog.DecrementStock(ctx, it.ProductID, it.Qty); derr != nil {
			s.restock(ctx, reserved)
			return Order{}, derr
		}
		reserved = append(reserved, it)
	}
	txID, payErr := s.payment.Charge(c.Subtotal, "IDR", failPayment)
	status := StatusPaid
	if payErr != nil {
		status = StatusPaymentFailed
		s.restock(ctx, reserved)
	}
	id, err := s.store.NextID(ctx)
	if err != nil {
		return Order{}, err
	}
	o := Order{
		ID:         id,
		SessionID:  sessionID,
		Email:      email,
		Items:      c.Items,
		TotalMinor: c.Subtotal,
		Currency:   "IDR",
		Status:     status,
		PaymentTx:  txID,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.store.Save(ctx, o); err != nil {
		return Order{}, err
	}
	if payErr != nil {
		return o, ErrChargeFail
	}
	if cerr := s.carts.Clear(ctx, sessionID); cerr != nil {
		return o, cerr
	}
	return o, nil
}

func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Order, error) {
	return s.store.List(ctx)
}

// SetStatus enforces a linear lifecycle.
func (s *Service) SetStatus(ctx context.Context, id string, next Status) (Order, error) {
	o, err := s.store.Get(ctx, id)
	if err != nil {
		return Order{}, err
	}
	allowed := map[Status][]Status{
		StatusPending:       {StatusPaid, StatusPaymentFailed, StatusCancelled},
		StatusPaid:          {StatusShipped, StatusCancelled},
		StatusPaymentFailed: {StatusCancelled},
		StatusShipped:       {StatusDone},
		StatusDone:          {},
		StatusCancelled:     {},
	}
	okTransition := false
	for _, n := range allowed[o.Status] {
		if n == next {
			okTransition = true
			break
		}
	}
	if !okTransition {
		return o, ErrBadStatus
	}
	// Cancelling a paid (unshipped) order releases the reservation.
	if o.Status == StatusPaid && next == StatusCancelled {
		s.restock(ctx, o.Items)
	}
	o.Status = next
	if err := s.store.Save(ctx, o); err != nil {
		return Order{}, err
	}
	return o, nil
}

// restock best-effort compensates reserved items; errors are swallowed
// because the order outcome is already decided (logged by the caller path).
func (s *Service) restock(ctx context.Context, items []cart.Item) {
	for _, it := range items {
		_ = s.catalog.IncrementStock(ctx, it.ProductID, it.Qty)
	}
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
