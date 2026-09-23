// Package catalog owns product domain rules and persistence behind a Store interface.
package catalog

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

// Product is the sellable unit. Prices are in minor units (e.g. sen for IDR).
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category,omitempty"`
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

// ListParams pages, filters and sorts the catalog. Page starts at 1.
type ListParams struct {
	Q        string
	Category string
	Sort     string // name_asc (default) | price_asc | price_desc | newest
	Page     int
	PerPage  int
}

// ListResult is a paginated product page.
type ListResult struct {
	Items   []Product `json:"items"`
	Total   int       `json:"total"`
	Page    int       `json:"page"`
	PerPage int       `json:"perPage"`
}

// List filters by query/category in the service so memory and Postgres
// behave identically; sorts and paginates the result.
func (s *Service) List(ctx context.Context, p ListParams) (ListResult, error) {
	all, err := s.store.List(ctx)
	if err != nil {
		return ListResult{}, err
	}
	q := strings.ToLower(strings.TrimSpace(p.Q))
	cat := strings.ToLower(strings.TrimSpace(p.Category))
	filtered := make([]Product, 0, len(all))
	for _, pr := range all {
		if q != "" && !strings.Contains(strings.ToLower(pr.Name), q) &&
			!strings.Contains(strings.ToLower(pr.Description), q) {
			continue
		}
		if cat != "" && strings.ToLower(pr.Category) != cat {
			continue
		}
		filtered = append(filtered, pr)
	}
	switch p.Sort {
	case "price_asc":
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].PriceMinor < filtered[j].PriceMinor })
	case "price_desc":
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].PriceMinor > filtered[j].PriceMinor })
	case "newest":
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.After(filtered[j].CreatedAt) })
	default:
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Name < filtered[j].Name })
	}
	page, perPage := p.Page, p.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 12
	}
	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}
	items := filtered[start:end]
	if items == nil {
		items = []Product{}
	}
	return ListResult{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

// Categories returns the distinct non-empty categories, sorted.
func (s *Service) Categories(ctx context.Context) ([]string, error) {
	all, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	out := []string{}
	for _, pr := range all {
		if pr.Category == "" || seen[pr.Category] {
			continue
		}
		seen[pr.Category] = true
		out = append(out, pr.Category)
	}
	sort.Strings(out)
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
