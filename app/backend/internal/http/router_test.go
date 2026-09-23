package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/open-commerce/backend/internal/auth"
	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/order"
	"github.com/open-commerce/backend/internal/payment"
)

func testRouter() http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	cat := catalog.NewService(catalog.NewMemoryStore())
	carts := cart.NewService(cat)
	pay := payment.NewStub()
	return NewRouter(Deps{Catalog: cat, Cart: carts, Order: order.NewService(carts, pay), Auth: auth.NewService(), Pay: pay}, logger)
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]bool
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil || !body["ok"] {
		t.Fatalf("unexpected body: %v", rec.Body.String())
	}
}

func TestListProductsSeeded(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var products []catalog.Product
	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(products) < 3 {
		t.Fatalf("expected >=3 seeded products, got %d", len(products))
	}
}

func TestAdminRequiresAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/admin/orders", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized && rec.Code != http.StatusForbidden {
		t.Fatalf("expected 401/403, got %d", rec.Code)
	}
}
