//go:build integration

// Integration test against a real Postgres. Run with:
//
//	DATABASE_URL=postgres://opencommerce:opencommerce@localhost:5432/opencommerce?sslmode=disable \
//	  go test -tags integration ./internal/db/
package db

import (
	"context"
	"testing"
	"time"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/order"
	"github.com/open-commerce/backend/internal/payment"

	"os"
)

func TestPostgresEndToEnd(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := Open(ctx, url)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer pool.Close()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for _, q := range []string{
		`TRUNCATE orders, carts, products`,
		`ALTER SEQUENCE order_seq RESTART WITH 1`,
	} {
		if _, err := pool.Exec(ctx, q); err != nil {
			t.Fatalf("truncate: %v", err)
		}
	}

	catStore := catalog.NewPGStore(pool)
	if err := catStore.SeedIfEmpty(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}
	cat := catalog.NewService(catStore)
	products, err := cat.List(ctx, "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(products) != 3 {
		t.Fatalf("expected 3 seeded products, got %d", len(products))
	}

	created, err := cat.Create(ctx, catalog.Product{Name: "PG Widget", PriceMinor: 2500000, Stock: 7})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := cat.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Stock != 7 {
		t.Fatalf("expected stock 7, got %d", got.Stock)
	}

	carts := cart.NewService(cart.NewPGStore(pool), cat)
	sess := "pg-sess-1"
	c, err := carts.AddItem(ctx, sess, created.ID, 2)
	if err != nil {
		t.Fatalf("add item: %v", err)
	}
	if c.Subtotal != 5000000 {
		t.Fatalf("expected 5000000, got %d", c.Subtotal)
	}
	// Persistence across service instances (proves DB, not memory).
	carts2 := cart.NewService(cart.NewPGStore(pool), cat)
	c2, err := carts2.Get(ctx, sess)
	if err != nil {
		t.Fatalf("re-get cart: %v", err)
	}
	if c2.Subtotal != 5000000 {
		t.Fatalf("cart not persisted, subtotal %d", c2.Subtotal)
	}

	orders := order.NewService(order.NewPGStore(pool), carts2, cat, payment.NewStub())
	o, err := orders.Checkout(ctx, sess, "pg@test.local", false)
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	if o.Status != order.StatusPaid || o.TotalMinor != 5000000 {
		t.Fatalf("unexpected order: %+v", o)
	}
	o, err = orders.SetStatus(ctx, o.ID, order.StatusShipped)
	if err != nil {
		t.Fatalf("ship: %v", err)
	}
	if o.Status != order.StatusShipped {
		t.Fatalf("expected shipped, got %s", o.Status)
	}
	list, err := orders.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 order, got %d (%v)", len(list), err)
	}
	got, err = cat.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("re-get product: %v", err)
	}
	if got.Stock != 5 {
		t.Fatalf("expected stock 5 after checkout, got %d", got.Stock)
	}

	// Oversell: stock 1, two sessions race — second checkout must fail.
	limited, err := cat.Create(ctx, catalog.Product{Name: "Limited", PriceMinor: 1000000, Stock: 1})
	if err != nil {
		t.Fatalf("create limited: %v", err)
	}
	if _, err := carts.AddItem(ctx, "pg-a", limited.ID, 1); err != nil {
		t.Fatalf("add a: %v", err)
	}
	if _, err := carts.AddItem(ctx, "pg-b", limited.ID, 1); err != nil {
		t.Fatalf("add b: %v", err)
	}
	if _, err := orders.Checkout(ctx, "pg-a", "a@test.local", false); err != nil {
		t.Fatalf("checkout a: %v", err)
	}
	if _, err := orders.Checkout(ctx, "pg-b", "b@test.local", false); err == nil {
		t.Fatal("expected oversell checkout to fail")
	}
}
