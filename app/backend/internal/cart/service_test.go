package cart

import (
	"context"
	"testing"

	"github.com/open-commerce/backend/internal/catalog"
)

func testCatalog(t *testing.T) *catalog.Service {
	t.Helper()
	c := catalog.NewService(catalog.NewMemoryStore())
	if _, err := c.Create(context.Background(), catalog.Product{ID: "p-test", Name: "Test", PriceMinor: 10000, Stock: 5}); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	return c
}

func TestAddItemSubtotal(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryStore(), testCatalog(t))
	c, err := svc.AddItem(ctx, "s-1", "p-test", 2)
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if c.Subtotal != 20000 {
		t.Fatalf("expected 20000, got %d", c.Subtotal)
	}
	c, err = svc.AddItem(ctx, "s-1", "p-test", 1)
	if err != nil {
		t.Fatal(err)
	}
	if c.Subtotal != 30000 {
		t.Fatalf("expected 30000 after increment, got %d", c.Subtotal)
	}
}

func TestAddItemExceedsStock(t *testing.T) {
	svc := NewService(NewMemoryStore(), testCatalog(t))
	if _, err := svc.AddItem(context.Background(), "s-1", "p-test", 99); err != ErrOutOfStock {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
}

func TestGetMissingIsEmpty(t *testing.T) {
	svc := NewService(NewMemoryStore(), testCatalog(t))
	c, err := svc.Get(context.Background(), "nope")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Items) != 0 || c.Subtotal != 0 {
		t.Fatalf("expected empty cart, got %+v", c)
	}
}

func TestSetQty(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryStore(), testCatalog(t))
	if _, err := svc.AddItem(ctx, "s-q", "p-test", 1); err != nil {
		t.Fatal(err)
	}
	c, err := svc.SetQty(ctx, "s-q", "p-test", 3)
	if err != nil {
		t.Fatal(err)
	}
	if c.Subtotal != 30000 {
		t.Fatalf("expected 30000, got %d", c.Subtotal)
	}
	c, err = svc.SetQty(ctx, "s-q", "p-test", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Items) != 0 {
		t.Fatalf("expected removed item, got %+v", c)
	}
	if _, err := svc.SetQty(ctx, "s-q", "p-missing", 1); err != ErrNoProduct {
		t.Fatalf("expected ErrNoProduct, got %v", err)
	}
	if _, err := svc.AddItem(ctx, "s-q", "p-test", 1); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SetQty(ctx, "s-q", "p-test", 99); err != ErrOutOfStock {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
}
