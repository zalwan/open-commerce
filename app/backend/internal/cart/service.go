// Package cart owns session carts. No stock reservation in v0.1.
package cart

import (
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
)

// Service validates against catalog and keeps carts in memory.
type Service struct {
	mu      sync.Mutex
	carts   map[string]*Cart
	catalog *catalog.Service
}

func NewService(catalogSvc *catalog.Service) *Service {
	return &Service{carts: make(map[string]*Cart), catalog: catalogSvc}
}

func (s *Service) Get(sessionID string) Cart {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.carts[sessionID]
	if !ok {
		return Cart{SessionID: sessionID, Items: []Item{}}
	}
	return *c
}

func (s *Service) AddItem(sessionID, productID string, qty int) (Cart, error) {
	if sessionID == "" {
		return Cart{}, ErrEmptySession
	}
	if qty <= 0 {
		return Cart{}, ErrBadQty
	}
	p, err := s.catalog.Get(productID)
	if err != nil {
		return Cart{}, ErrNoProduct
	}
	if qty > p.Stock {
		return Cart{}, ErrOutOfStock
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.carts[sessionID]
	if !ok {
		c = &Cart{SessionID: sessionID}
		s.carts[sessionID] = c
	}
	for i, it := range c.Items {
		if it.ProductID == productID {
			c.Items[i].Qty += qty
			c.Items[i].PriceMinor = p.PriceMinor
			c.Items[i].Name = p.Name
			c.recalc()
			return *c, nil
		}
	}
	c.Items = append(c.Items, Item{ProductID: p.ID, Name: p.Name, PriceMinor: p.PriceMinor, Qty: qty})
	c.recalc()
	return *c, nil
}

func (s *Service) RemoveItem(sessionID, productID string) Cart {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.carts[sessionID]
	if !ok {
		return Cart{SessionID: sessionID, Items: []Item{}}
	}
	kept := c.Items[:0]
	for _, it := range c.Items {
		if it.ProductID != productID {
			kept = append(kept, it)
		}
	}
	c.Items = kept
	c.recalc()
	return *c
}

func (s *Service) Clear(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.carts, sessionID)
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
