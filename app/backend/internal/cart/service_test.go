package cart

import (
	"testing"

	"github.com/open-commerce/backend/internal/catalog"
)

func testCatalog(t *testing.T) *catalog.Service {
	t.Helper()
	c := catalog.NewService(catalog.NewMemoryStore())
	if _, err := c.Create(catalog.Product{ID: "p-test", Name: "Test", PriceMinor: 10000, Stock: 5}); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	return c
}

func TestAddItemSubtotal(t *testing.T) {
	svc := NewService(testCatalog(t))
	c, err := svc.AddItem("s-1", "p-test", 2)
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if c.Subtotal != 20000 {
		t.Fatalf("expected 20000, got %d", c.Subtotal)
	}
	c, _ = svc.AddItem("s-1", "p-test", 1)
	if c.Subtotal != 30000 {
		t.Fatalf("expected 30000 after increment, got %d", c.Subtotal)
	}
}

func TestAddItemExceedsStock(t *testing.T) {
	svc := NewService(testCatalog(t))
	if _, err := svc.AddItem("s-1", "p-test", 99); err != ErrOutOfStock {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
}
