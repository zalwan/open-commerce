package catalog

import (
	"context"
	"testing"
)

func TestValidate(t *testing.T) {
	if err := Validate(Product{Name: "A", PriceMinor: 100, Stock: 1}); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	if err := Validate(Product{PriceMinor: 100}); err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
	if err := Validate(Product{Name: "A", PriceMinor: 0}); err != ErrInvalidPrice {
		t.Fatalf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestCreateAndList(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryStore())
	before, err := svc.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	p, err := svc.Create(ctx, Product{Name: "Keyboard", PriceMinor: 5000000, Stock: 5})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if p.ID == "" || p.Currency != "IDR" {
		t.Fatalf("unexpected product: %+v", p)
	}
	hits, err := svc.List(ctx, "keyboard")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 search hit, got %d", len(hits))
	}
	all, err := svc.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != len(before)+1 {
		t.Fatalf("expected %d products, got %d", len(before)+1, len(all))
	}
}

func TestCreateDuplicate(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryStore())
	if _, err := svc.Create(ctx, Product{ID: "p-dup", Name: "Dup", PriceMinor: 100, Stock: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Create(ctx, Product{ID: "p-dup", Name: "Dup", PriceMinor: 100, Stock: 1}); err != ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}
