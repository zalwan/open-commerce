package order

import (
	"testing"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/payment"
)

func setup(t *testing.T) (*Service, string) {
	t.Helper()
	cat := catalog.NewService(catalog.NewMemoryStore())
	if _, err := cat.Create(catalog.Product{ID: "p-a", Name: "A", PriceMinor: 10000, Stock: 10}); err != nil {
		t.Fatal(err)
	}
	carts := cart.NewService(cat)
	if _, err := carts.AddItem("sess-1", "p-a", 2); err != nil {
		t.Fatal(err)
	}
	return NewService(carts, payment.NewStub()), "sess-1"
}

func TestCheckoutSuccessClearsCart(t *testing.T) {
	svc, sess := setup(t)
	o, err := svc.Checkout(sess, "buyer@mail.test", false)
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}
	if o.Status != StatusPaid || o.TotalMinor != 20000 || o.PaymentTx == "" {
		t.Fatalf("unexpected order: %+v", o)
	}
}

func TestCheckoutPaymentFailKeepsOrder(t *testing.T) {
	svc, sess := setup(t)
	o, err := svc.Checkout(sess, "buyer@mail.test", true)
	if err != ErrChargeFail {
		t.Fatalf("expected ErrChargeFail, got %v", err)
	}
	if o.Status != StatusPaymentFailed {
		t.Fatalf("expected payment_failed, got %s", o.Status)
	}
}
