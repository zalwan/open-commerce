package httpapi

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-commerce/backend/internal/auth"
	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/media"
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
	var res catalog.ListResult
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(res.Items) < 3 || res.Total < 3 {
		t.Fatalf("expected >=3 seeded products, got %+v", res)
	}
	// Category filter + categories endpoint.
	req = httptest.NewRequest("GET", "/api/v1/products?category=Fashion", nil)
	rec = httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	var filtered catalog.ListResult
	if err := json.NewDecoder(rec.Body).Decode(&filtered); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if filtered.Total != 2 {
		t.Fatalf("expected 2 fashion, got %+v", filtered)
	}
	req = httptest.NewRequest("GET", "/api/v1/categories", nil)
	rec = httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	var cats []string
	if err := json.NewDecoder(rec.Body).Decode(&cats); err != nil || len(cats) != 2 {
		t.Fatalf("expected 2 categories, got %v (%v)", cats, err)
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

// tiny 1x1 PNG for upload tests.
var testPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00,
	0x0c, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xff, 0xff, 0x3f,
	0x00, 0x05, 0xfe, 0x02, 0xfe, 0xdc, 0xcc, 0x59, 0xe7, 0x00, 0x00, 0x00,
	0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

func testRouterWithMedia(t *testing.T) (http.Handler, string) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	cat := catalog.NewService(catalog.NewMemoryStore())
	carts := cart.NewService(cart.NewMemoryStore(), cat)
	pay := payment.NewStub()
	dir := t.TempDir()
	return NewRouter(Deps{
		Catalog: cat, Cart: carts,
		Order: order.NewService(order.NewMemoryStore(), carts, cat, pay),
		Auth:  auth.NewService(), Pay: pay,
		Media: media.NewStore(dir, "/static/"),
	}, logger), dir
}

func uploadReq(t *testing.T, url, token, field, filename string, content []byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile(field, filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", url, &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestUploadImage(t *testing.T) {
	router, dir := testRouterWithMedia(t)

	loginBody := bytes.NewBufferString(`{"email":"admin@shop.test","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", loginBody)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var login struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&login); err != nil {
		t.Fatal(err)
	}

	// Unknown product -> 404 (file still stored, harmless temp).
	req = uploadReq(t, "/api/v1/admin/products/nope/image", login.Token, "image", "a.png", testPNG)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	// Create product then upload.
	createBody := bytes.NewBufferString(`{"id":"p-img","name":"Img","priceMinor":1000,"stock":2}`)
	req = httptest.NewRequest("POST", "/api/v1/admin/products", createBody)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create failed: %d", rec.Code)
	}
	req = uploadReq(t, "/api/v1/admin/products/p-img/image", login.Token, "image", "a.png", testPNG)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("upload failed: %d %s", rec.Code, rec.Body.String())
	}
	var p catalog.Product
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil || p.ImageURL != "/static/p-img.png" {
		t.Fatalf("unexpected product: %+v (%v)", p, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "p-img.png")); err != nil {
		t.Fatalf("file not stored: %v", err)
	}

	// Serve it back.
	req = httptest.NewRequest("GET", "/static/p-img.png", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("static serve failed: %d %s", rec.Code, rec.Header().Get("Content-Type"))
	}

	// Bad type -> 400.
	req = uploadReq(t, "/api/v1/admin/products/p-img/image", login.Token, "image", "a.txt", []byte("hello"))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	// No token -> 401/403.
	req = uploadReq(t, "/api/v1/admin/products/p-img/image", "", "image", "a.png", testPNG)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized && rec.Code != http.StatusForbidden {
		t.Fatalf("expected 401/403, got %d", rec.Code)
	}
}

func TestAdminStats(t *testing.T) {
	router := testRouter()
	loginBody := bytes.NewBufferString(`{"email":"admin@shop.test","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", loginBody)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var login struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&login); err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var stats map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if stats["products"].(float64) < 3 {
		t.Fatalf("expected >=3 products, got %v", stats)
	}
}
