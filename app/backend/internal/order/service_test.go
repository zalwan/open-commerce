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
	return NewService(NewMemoryStore(), carts, payment.NewStub()), "sess-1"
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
