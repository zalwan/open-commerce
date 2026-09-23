// Package catalog memory store: ephemeral v0.1 persistence with seed data.
package catalog

import (
	"sort"
	"sync"
	"time"
)

// MemoryStore is a mutex-guarded map. Replace with Postgres in v0.2.
type MemoryStore struct {
	mu sync.RWMutex
	m  map[string]Product
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{m: make(map[string]Product)}
	s.seed()
	return s
}

func (s *MemoryStore) List() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Product, 0, len(s.m))
	for _, p := range s.m {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (s *MemoryStore) Get(id string) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.m[id]
	return p, ok
}

func (s *MemoryStore) Save(p Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[p.ID] = p
}

func (s *MemoryStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[id]; !ok {
		return false
	}
	delete(s.m, id)
	return true
}

func (s *MemoryStore) seed() {
	now := time.Now().UTC()
	for _, p := range []Product{
		{ID: "p-kaos-hitam", Name: "Kaos Hitam", Description: "Kaos katun premium", PriceMinor: 12900000, Currency: "IDR", Stock: 50, CreatedAt: now},
		{ID: "p-kopi-arabika", Name: "Kopi Arabika 250g", Description: "Biji kopi sangrai medium", PriceMinor: 8500000, Currency: "IDR", Stock: 100, CreatedAt: now},
		{ID: "p-tas-kanvas", Name: "Tas Kanvas", Description: "Tas selempang kanvas", PriceMinor: 19900000, Currency: "IDR", Stock: 20, CreatedAt: now},
	} {
		s.m[p.ID] = p
	}
}
