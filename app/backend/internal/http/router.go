// Package httpapi wires services to REST routes. Handlers never touch stores directly.
package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/open-commerce/backend/internal/auth"
	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/media"
	"github.com/open-commerce/backend/internal/order"
	"github.com/open-commerce/backend/internal/payment"
)

// Deps are wired in main.go.
type Deps struct {
	Catalog *catalog.Service
	Cart    *cart.Service
	Order   *order.Service
	Auth    *auth.Service
	Pay     payment.Provider
	Media   *media.Store
	// Ready is nil on memory store; set to a DB ping when Postgres is configured.
	Ready func(ctx context.Context) error
}

// NewRouter builds the mux with logging + CORS middleware.
func NewRouter(d Deps, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	h := &handlers{deps: d}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
	mux.HandleFunc("GET /readyz", h.readyz)
	// Catalog (public read).
	mux.HandleFunc("GET /api/v1/products", h.listProducts)
	mux.HandleFunc("GET /api/v1/categories", h.listCategories)
	mux.HandleFunc("GET /api/v1/products/{id}", h.getProduct)
	// Cart.
	mux.HandleFunc("GET /api/v1/cart", h.getCart)
	mux.HandleFunc("POST /api/v1/cart/items", h.addCartItem)
	mux.HandleFunc("PUT /api/v1/cart/items", h.setCartItemQty)
	mux.HandleFunc("DELETE /api/v1/cart/items/{productID}", h.removeCartItem)
	// Orders.
	mux.HandleFunc("POST /api/v1/orders/checkout", h.checkout)
	mux.HandleFunc("GET /api/v1/orders", h.listMyOrders)
	mux.HandleFunc("GET /api/v1/orders/{id}", h.getOrder)
	// Auth stub.
	mux.HandleFunc("POST /api/v1/auth/login", h.login)
	mux.HandleFunc("GET /api/v1/auth/me", h.me)
	// Admin (auth stub gated).
	mux.HandleFunc("POST /api/v1/admin/products", h.requireAdmin(h.createProduct))
	mux.HandleFunc("PUT /api/v1/admin/products/{id}", h.requireAdmin(h.updateProduct))
	mux.HandleFunc("DELETE /api/v1/admin/products/{id}", h.requireAdmin(h.deleteProduct))
	mux.HandleFunc("POST /api/v1/admin/products/{id}/image", h.requireAdmin(h.uploadImage))
	mux.HandleFunc("GET /api/v1/admin/orders", h.requireAdmin(h.listOrders))
	mux.HandleFunc("GET /api/v1/admin/stats", h.requireAdmin(h.adminStats))
	mux.HandleFunc("POST /api/v1/admin/orders/{id}/status", h.requireAdmin(h.setOrderStatus))
	// Public static files (product uploads). No directory listing.
	if d.Media != nil {
		mux.HandleFunc("GET /static/", h.serveStatic)
	}

	return withLogging(withCORS(mux), logger)
}

type handlers struct {
	deps Deps
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request", "method", r.Method, "path", r.URL.Path, "latency_ms", time.Since(start).Milliseconds())
	})
}

// requireAdmin enforces Bearer auth-stub with admin role.
func (h *handlers) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			writeErr(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		sess, err := h.deps.Auth.Authenticate(token)
		if err != nil || sess.Role != auth.RoleAdmin {
			writeErr(w, http.StatusForbidden, "admin required")
			return
		}
		next(w, r)
	}
}
