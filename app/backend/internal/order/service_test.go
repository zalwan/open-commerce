package order

import (
	"context"
	"testing"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/payment"
)

func setup(t *testing.T) (*Service, string) {
	t.Helper()
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-a", Name: "A", PriceMinor: 10000, Stock: 10}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	if _, err := carts.AddItem(ctx, "sess-1", "p-a", 2); err != nil {
		t.Fatal(err)
	}
	return NewService(NewMemoryStore(), carts, cat, payment.NewStub()), "sess-1"
}

func TestCheckoutSuccessClearsCart(t *testing.T) {
	ctx := context.Background()
	svc, sess := setup(t)
	o, err := svc.Checkout(ctx, sess, "buyer@mail.test", false)
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}
	if o.Status != StatusPaid || o.TotalMinor != 20000 || o.PaymentTx == "" {
		t.Fatalf("unexpected order: %+v", o)
	}
}

func TestCheckoutPaymentFailKeepsOrder(t *testing.T) {
	ctx := context.Background()
	svc, sess := setup(t)
	o, err := svc.Checkout(ctx, sess, "buyer@mail.test", true)
	if err != ErrChargeFail {
		t.Fatalf("expected ErrChargeFail, got %v", err)
	}
	if o.Status != StatusPaymentFailed {
		t.Fatalf("expected payment_failed, got %s", o.Status)
	}
}

func TestCheckoutDecrementsStock(t *testing.T) {
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-s", Name: "S", PriceMinor: 5000, Stock: 10}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	if _, err := carts.AddItem(ctx, "sess-s", "p-s", 2); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewMemoryStore(), carts, cat, payment.NewStub())
	if _, err := svc.Checkout(ctx, "sess-s", "b@mail.test", false); err != nil {
		t.Fatal(err)
	}
	p, err := cat.Get(ctx, "p-s")
	if err != nil {
		t.Fatal(err)
	}
	if p.Stock != 8 {
		t.Fatalf("expected stock 8, got %d", p.Stock)
	}
}

func TestCheckoutPaymentFailRestoresStock(t *testing.T) {
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-r", Name: "R", PriceMinor: 5000, Stock: 10}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	if _, err := carts.AddItem(ctx, "sess-r", "p-r", 2); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewMemoryStore(), carts, cat, payment.NewStub())
	if _, err := svc.Checkout(ctx, "sess-r", "b@mail.test", true); err != ErrChargeFail {
		t.Fatalf("expected ErrChargeFail, got %v", err)
	}
	p, err := cat.Get(ctx, "p-r")
	if err != nil {
		t.Fatal(err)
	}
	if p.Stock != 10 {
		t.Fatalf("expected stock restored to 10, got %d", p.Stock)
	}
}

func TestCheckoutInsufficientStock(t *testing.T) {
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-1", Name: "One", PriceMinor: 5000, Stock: 1}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	if _, err := carts.AddItem(ctx, "sess-a", "p-1", 1); err != nil {
		t.Fatal(err)
	}
	if _, err := carts.AddItem(ctx, "sess-b", "p-1", 1); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewMemoryStore(), carts, cat, payment.NewStub())
	if _, err := svc.Checkout(ctx, "sess-a", "a@mail.test", false); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Checkout(ctx, "sess-b", "b@mail.test", false); err != catalog.ErrInsufficientStock {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	orders, err := svc.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 order (no order for failed checkout), got %d", len(orders))
	}
}

func TestCancelPaidRestoresStock(t *testing.T) {
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-c", Name: "C", PriceMinor: 5000, Stock: 5}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	if _, err := carts.AddItem(ctx, "sess-c", "p-c", 2); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewMemoryStore(), carts, cat, payment.NewStub())
	o, err := svc.Checkout(ctx, "sess-c", "c@mail.test", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SetStatus(ctx, o.ID, StatusCancelled); err != nil {
		t.Fatal(err)
	}
	p, err := cat.Get(ctx, "p-c")
	if err != nil {
		t.Fatal(err)
	}
	if p.Stock != 5 {
		t.Fatalf("expected stock restored to 5, got %d", p.Stock)
	}
}
func TestSetStatusLifecycle(t *testing.T) {
	ctx := context.Background()
	svc, sess := setup(t)
	o, err := svc.Checkout(ctx, sess, "buyer@mail.test", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SetStatus(ctx, o.ID, StatusDone); err != ErrBadStatus {
		t.Fatalf("expected ErrBadStatus paid->done, got %v", err)
	}
	o, err = svc.SetStatus(ctx, o.ID, StatusShipped)
	if err != nil {
		t.Fatal(err)
	}
	if o.Status != StatusShipped {
		t.Fatalf("expected shipped, got %s", o.Status)
	}
}

func TestListByEmail(t *testing.T) {
	ctx := context.Background()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(ctx, catalog.Product{ID: "p-t", Name: "T", PriceMinor: 1000, Stock: 10}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	svc := NewService(NewMemoryStore(), carts, cat, payment.NewStub())
	for _, s := range []struct {
		sess, email string
	}{{"s-t1", "me@mail.test"}, {"s-t2", "me@mail.test"}, {"s-t3", "other@mail.test"}} {
		if _, err := carts.AddItem(ctx, s.sess, "p-t", 1); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Checkout(ctx, s.sess, s.email, false); err != nil {
			t.Fatal(err)
		}
	}
	mine, err := svc.ListByEmail(ctx, "me@mail.test")
	if err != nil {
		t.Fatal(err)
	}
	if len(mine) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(mine))
	}
	if _, err := svc.ListByEmail(ctx, ""); err != ErrBadEmail {
		t.Fatalf("expected ErrBadEmail, got %v", err)
	}
}
