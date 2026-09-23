// Package catalog owns product domain rules and persistence behind a Store interface.
package catalog

import (
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

// Store abstracts persistence; v0.1 uses memory, v0.2+ Postgres.
type Store interface {
	List() []Product
	Get(id string) (Product, bool)
	Save(p Product)
	Delete(id string) bool
}

var (
	ErrNotFound      = errors.New("product not found")
	ErrInvalidName   = errors.New("name is required")
	ErrInvalidPrice  = errors.New("priceMinor must be > 0")
	ErrInvalidStock  = errors.New("stock must be >= 0")
	ErrAlreadyExists = errors.New("product already exists")
)

// Service enforces catalog rules; handlers must go through it.
type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

// List returns all products, optionally filtered by case-insensitive name query.
func (s *Service) List(q string) []Product {
	all := s.store.List()
	if q == "" {
		return all
	}
	q = strings.ToLower(q)
	out := make([]Product, 0, len(all))
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Name), q) {
			out = append(out, p)
		}
	}
	return out
}

func (s *Service) Get(id string) (Product, error) {
	p, ok := s.store.Get(id)
	if !ok {
		return Product{}, ErrNotFound
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

func (s *Service) Create(p Product) (Product, error) {
	if err := Validate(p); err != nil {
		return Product{}, err
	}
	if _, ok := s.store.Get(p.ID); ok && p.ID != "" {
		return Product{}, ErrAlreadyExists
	}
	if p.ID == "" {
		p.ID = "p-" + strings.ToLower(strings.ReplaceAll(strings.TrimSpace(p.Name), " ", "-"))
	}
	if p.Currency == "" {
		p.Currency = "IDR"
	}
	p.CreatedAt = time.Now().UTC()
	s.store.Save(p)
	return p, nil
}

func (s *Service) Update(p Product) (Product, error) {
	if err := Validate(p); err != nil {
		return Product{}, err
	}
	existing, ok := s.store.Get(p.ID)
	if !ok {
		return Product{}, ErrNotFound
	}
	p.CreatedAt = existing.CreatedAt
	if p.Currency == "" {
		p.Currency = existing.Currency
	}
	s.store.Save(p)
	return p, nil
}

func (s *Service) Delete(id string) error {
	if !s.store.Delete(id) {
		return ErrNotFound
	}
	return nil
}
