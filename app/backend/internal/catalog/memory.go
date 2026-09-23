// Memory store: ephemeral v0.1 persistence with seed data.
package catalog

import (
	"context"
	"sort"
	"sync"
	"time"
)

// MemoryStore is a mutex-guarded map. Used when DATABASE_URL is empty.
type MemoryStore struct {
	mu sync.RWMutex
	m  map[string]Product
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{m: make(map[string]Product)}
	s.seed()
	return s
}

func (s *MemoryStore) List(_ context.Context) ([]Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Product, 0, len(s.m))
	for _, p := range s.m {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (s *MemoryStore) Get(_ context.Context, id string) (Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) Save(_ context.Context, p Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[p.ID] = p
	return nil
}

func (s *MemoryStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

// SeedProducts is shared with PGStore.SeedIfEmpty so both backends start identical.
func SeedProducts(now time.Time) []Product {
	return []Product{
		{ID: "p-kaos-hitam", Name: "Kaos Hitam", Description: "Kaos katun premium", PriceMinor: 12900000, Currency: "IDR", Stock: 50, CreatedAt: now},
		{ID: "p-kopi-arabika", Name: "Kopi Arabika 250g", Description: "Biji kopi sangrai medium", PriceMinor: 8500000, Currency: "IDR", Stock: 100, CreatedAt: now},
		{ID: "p-tas-kanvas", Name: "Tas Kanvas", Description: "Tas selempang kanvas", PriceMinor: 19900000, Currency: "IDR", Stock: 20, CreatedAt: now},
	}
}

func (s *MemoryStore) seed() {
	for _, p := range SeedProducts(time.Now().UTC()) {
		s.m[p.ID] = p
	}
}
