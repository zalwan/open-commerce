// Package cart owns session carts. No stock reservation in v0.2.
package cart

import (
	"context"
	"errors"
	"sync"

	"github.com/open-commerce/backend/internal/catalog"
)

// Item is a price snapshot at add-time.
type Item struct {
	ProductID  string `json:"productId"`
	Name       string `json:"name"`
	PriceMinor int64  `json:"priceMinor"`
	Qty        int    `json:"qty"`
}

// Cart is one session's basket.
type Cart struct {
	SessionID string `json:"sessionId"`
	Items     []Item `json:"items"`
	Subtotal  int64  `json:"subtotalMinor"`
}

var (
	ErrEmptySession = errors.New("session_id is required")
	ErrBadQty       = errors.New("qty must be > 0")
	ErrNoProduct    = errors.New("product not found")
	ErrOutOfStock   = errors.New("qty exceeds stock")
	ErrCartNotFound = errors.New("cart not found")
)

// Store abstracts cart persistence: memory (dev) or Postgres (DATABASE_URL set).
type Store interface {
	Get(ctx context.Context, sessionID string) (Cart, error)
	Save(ctx context.Context, c Cart) error
	Delete(ctx context.Context, sessionID string) error
}

// MemoryStore keeps carts process-local. Used when DATABASE_URL is empty.
type MemoryStore struct {
	mu    sync.Mutex
	carts map[string]*Cart
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{carts: make(map[string]*Cart)} }

func (s *MemoryStore) Get(_ context.Context, sessionID string) (Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.carts[sessionID]
	if !ok {
		return Cart{}, ErrCartNotFound
	}
	cp := *c
	return cp, nil
}

func (s *MemoryStore) Save(_ context.Context, c Cart) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := c
	s.carts[c.SessionID] = &cp
	return nil
}

func (s *MemoryStore) Delete(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.carts, sessionID)
	return nil
}

// Service validates against catalog and delegates persistence to a Store.
type Service struct {
	store   Store
	catalog *catalog.Service
}

func NewService(store Store, catalogSvc *catalog.Service) *Service {
	return &Service{store: store, catalog: catalogSvc}
}

func (s *Service) Get(ctx context.Context, sessionID string) (Cart, error) {
	if sessionID == "" {
		return Cart{}, ErrEmptySession
	}
	c, err := s.store.Get(ctx, sessionID)
	if errors.Is(err, ErrCartNotFound) {
		return Cart{SessionID: sessionID, Items: []Item{}}, nil
	}
	return c, err
}

func (s *Service) AddItem(ctx context.Context, sessionID, productID string, qty int) (Cart, error) {
	if sessionID == "" {
		return Cart{}, ErrEmptySession
	}
	if qty <= 0 {
		return Cart{}, ErrBadQty
	}
	p, err := s.catalog.Get(ctx, productID)
	if err != nil {
		return Cart{}, ErrNoProduct
	}
	if qty > p.Stock {
		return Cart{}, ErrOutOfStock
	}
	c, err := s.Get(ctx, sessionID)
	if err != nil {
		return Cart{}, err
	}
	merged := false
	for i, it := range c.Items {
		if it.ProductID == productID {
			c.Items[i].Qty += qty
			c.Items[i].PriceMinor = p.PriceMinor
			c.Items[i].Name = p.Name
			merged = true
			break
		}
	}
	if !merged {
		c.Items = append(c.Items, Item{ProductID: p.ID, Name: p.Name, PriceMinor: p.PriceMinor, Qty: qty})
	}
	c.recalc()
	if err := s.store.Save(ctx, c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (s *Service) RemoveItem(ctx context.Context, sessionID, productID string) (Cart, error) {
	c, err := s.Get(ctx, sessionID)
	if err != nil {
		return Cart{}, err
	}
	kept := c.Items[:0]
	for _, it := range c.Items {
		if it.ProductID != productID {
			kept = append(kept, it)
		}
	}
	c.Items = kept
	c.recalc()
	if err := s.store.Save(ctx, c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (s *Service) Clear(ctx context.Context, sessionID string) error {
	return s.store.Delete(ctx, sessionID)
}

// SetQty sets an absolute quantity; qty 0 removes the item.
// The item must already be in the cart, else ErrNoProduct.
func (s *Service) SetQty(ctx context.Context, sessionID, productID string, qty int) (Cart, error) {
	if qty < 0 {
		return Cart{}, ErrBadQty
	}
	if qty == 0 {
		return s.RemoveItem(ctx, sessionID, productID)
	}
	p, err := s.catalog.Get(ctx, productID)
	if err != nil {
		return Cart{}, ErrNoProduct
	}
	if qty > p.Stock {
		return Cart{}, ErrOutOfStock
	}
	c, err := s.Get(ctx, sessionID)
	if err != nil {
		return Cart{}, err
	}
	found := false
	for i, it := range c.Items {
		if it.ProductID == productID {
			c.Items[i].Qty = qty
			c.Items[i].PriceMinor = p.PriceMinor
			c.Items[i].Name = p.Name
			found = true
			break
		}
	}
	if !found {
		return Cart{}, ErrNoProduct
	}
	c.recalc()
	if err := s.store.Save(ctx, c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (c *Cart) recalc() {
	var sum int64
	for _, it := range c.Items {
		sum += it.PriceMinor * int64(it.Qty)
	}
	c.Subtotal = sum
	if c.Items == nil {
		c.Items = []Item{}
	}
}
