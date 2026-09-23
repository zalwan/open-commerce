// Package catalog owns product domain rules and persistence behind a Store interface.
package catalog

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Product is the sellable unit. Prices are in minor units (e.g. sen for IDR).
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	PriceMinor  int64     `json:"priceMinor"`
	Currency    string    `json:"currency"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Store abstracts persistence: MemoryStore (dev/test) or PGStore (DATABASE_URL set).
type Store interface {
	List(ctx context.Context) ([]Product, error)
	Get(ctx context.Context, id string) (Product, error)
	Save(ctx context.Context, p Product) error
	Delete(ctx context.Context, id string) error
	// DecrementStock atomically reduces stock, failing with
	// ErrInsufficientStock when stock < qty. Used by checkout.
	DecrementStock(ctx context.Context, id string, qty int) error
	// IncrementStock restores stock (payment failure compensation, cancel).
	IncrementStock(ctx context.Context, id string, qty int) error
}

var (
	ErrNotFound          = errors.New("product not found")
	ErrInvalidName       = errors.New("name is required")
	ErrInvalidPrice      = errors.New("priceMinor must be > 0")
	ErrInvalidStock      = errors.New("stock must be >= 0")
	ErrAlreadyExists     = errors.New("product already exists")
	ErrBadQty            = errors.New("qty must be > 0")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// Service enforces catalog rules; handlers must go through it.
type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

// List returns all products, optionally filtered by case-insensitive name query.
func (s *Service) List(ctx context.Context, q string) ([]Product, error) {
	all, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	if q == "" {
		return all, nil
	}
	q = strings.ToLower(q)
	out := make([]Product, 0, len(all))
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Name), q) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, id string) (Product, error) {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

func Validate(p Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrInvalidName
	}
	if p.PriceMinor <= 0 {
		return ErrInvalidPrice
	}
	if p.Stock < 0 {
		return ErrInvalidStock
	}
	return nil
}

func (s *Service) Create(ctx context.Context, p Product) (Product, error) {
	if err := Validate(p); err != nil {
		return Product{}, err
	}
	if p.ID != "" {
		if _, err := s.store.Get(ctx, p.ID); err == nil {
			return Product{}, ErrAlreadyExists
		}
	}
	if p.ID == "" {
		p.ID = "p-" + strings.ToLower(strings.ReplaceAll(strings.TrimSpace(p.Name), " ", "-"))
	}
	if p.Currency == "" {
		p.Currency = "IDR"
	}
	p.CreatedAt = time.Now().UTC()
	if err := s.store.Save(ctx, p); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, p Product) (Product, error) {
	if err := Validate(p); err != nil {
		return Product{}, err
	}
	existing, err := s.store.Get(ctx, p.ID)
	if err != nil {
		return Product{}, err
	}
	p.CreatedAt = existing.CreatedAt
	if p.Currency == "" {
		p.Currency = existing.Currency
	}
	if err := s.store.Save(ctx, p); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

// DecrementStock reserves stock for checkout. qty must be > 0.
func (s *Service) DecrementStock(ctx context.Context, id string, qty int) error {
	if qty <= 0 {
		return ErrBadQty
	}
	return s.store.DecrementStock(ctx, id, qty)
}

// IncrementStock restores previously reserved stock.
func (s *Service) IncrementStock(ctx context.Context, id string, qty int) error {
	if qty <= 0 {
		return ErrBadQty
	}
	return s.store.IncrementStock(ctx, id, qty)
}
