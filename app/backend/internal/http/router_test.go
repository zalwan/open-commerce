package httpapi

import (
	"bytes"
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
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	pay := payment.NewStub()
	return NewRouter(Deps{Catalog: cat, Cart: carts, Order: order.NewService(order.NewMemoryStore(), carts, cat, pay), Auth: auth.NewService(), Pay: pay}, logger)
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

func TestReadyzMemory(t *testing.T) {
	req := httptest.NewRequest("GET", "/readyz", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
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

func TestAdminProductCRUD(t *testing.T) {
	router := testRouter()

	// Login as admin.
	loginBody := bytes.NewBufferString(`{"email":"admin@shop.test","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", loginBody)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login failed: %d", rec.Code)
	}
	var login struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&login); err != nil || login.Token == "" {
		t.Fatalf("bad login response: %v", rec.Body.String())
	}
	authHeader := "Bearer " + login.Token

	// Create.
	req = httptest.NewRequest("POST", "/api/v1/admin/products", bytes.NewBufferString(`{"id":"p-test-crud","name":"CRUD","priceMinor":1000,"stock":3}`))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", rec.Code, rec.Body.String())
	}

	// Invalid create (price 0) -> 400.
	req = httptest.NewRequest("POST", "/api/v1/admin/products", bytes.NewBufferString(`{"id":"p-bad","name":"Bad","priceMinor":0,"stock":1}`))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// Delete.
	req = httptest.NewRequest("DELETE", "/api/v1/admin/products/p-test-crud", nil)
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete failed: %d %s", rec.Code, rec.Body.String())
	}
}
