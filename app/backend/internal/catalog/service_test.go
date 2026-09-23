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
	before, err := svc.List(ctx, ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	p, err := svc.Create(ctx, Product{Name: "Keyboard", PriceMinor: 5000000, Stock: 5, Category: "Gadgets"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if p.ID == "" || p.Currency != "IDR" {
		t.Fatalf("unexpected product: %+v", p)
	}
	hits, err := svc.List(ctx, ListParams{Q: "keyboard"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits.Items) != 1 {
		t.Fatalf("expected 1 search hit, got %d", len(hits.Items))
	}
	all, err := svc.List(ctx, ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Items) != len(before.Items)+1 {
		t.Fatalf("expected %d products, got %d", len(before.Items)+1, len(all.Items))
	}
}

func TestListCategorySortPaginate(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryStore())
	// 3 seeds: 2 Fashion + 1 Food & Drink.
	fashion, err := svc.List(ctx, ListParams{Category: "fashion", Sort: "price_desc"})
	if err != nil {
		t.Fatal(err)
	}
	if fashion.Total != 2 || len(fashion.Items) != 2 {
		t.Fatalf("expected 2 fashion, got %+v", fashion)
	}
	if fashion.Items[0].PriceMinor < fashion.Items[1].PriceMinor {
		t.Fatal("expected price_desc order")
	}
	page1, err := svc.List(ctx, ListParams{PerPage: 2, Page: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page1.Total != 3 || len(page1.Items) != 2 || page1.Page != 1 {
		t.Fatalf("unexpected page1: %+v", page1)
	}
	page2, err := svc.List(ctx, ListParams{PerPage: 2, Page: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Items) != 1 {
		t.Fatalf("expected 1 item on page 2, got %+v", page2)
	}
	cats, err := svc.Categories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 2 || cats[0] != "Fashion" {
		t.Fatalf("unexpected categories: %v", cats)
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
