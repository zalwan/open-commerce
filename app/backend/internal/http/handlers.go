package httpapi

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/media"
	"github.com/open-commerce/backend/internal/order"
)

// Domain errors map to 4xx; anything else is a 500 without leaking internals.
func catalogErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound):
		writeErr(w, http.StatusNotFound, "product not found")
	case errors.Is(err, catalog.ErrInvalidName),
		errors.Is(err, catalog.ErrInvalidPrice),
		errors.Is(err, catalog.ErrInvalidStock),
		errors.Is(err, catalog.ErrAlreadyExists):
		writeErr(w, http.StatusBadRequest, err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal error")
	}
}

func cartErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, cart.ErrEmptySession),
		errors.Is(err, cart.ErrBadQty),
		errors.Is(err, cart.ErrNoProduct),
		errors.Is(err, cart.ErrOutOfStock):
		writeErr(w, http.StatusBadRequest, err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal error")
	}
}

func orderErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, order.ErrNotFound):
		writeErr(w, http.StatusNotFound, "order not found")
	case errors.Is(err, order.ErrBadEmail),
		errors.Is(err, order.ErrEmptyCart),
		errors.Is(err, order.ErrBadStatus):
		writeErr(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, catalog.ErrInsufficientStock),
		errors.Is(err, catalog.ErrNotFound):
		writeErr(w, http.StatusConflict, err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal error")
	}
}

func (h *handlers) listProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	products, err := h.deps.Catalog.List(r.Context(), catalog.ListParams{
		Q:        q.Get("q"),
		Category: q.Get("category"),
		Sort:     q.Get("sort"),
		Page:     page,
		PerPage:  perPage,
	})
	if err != nil {
		catalogErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *handlers) listCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.deps.Catalog.Categories(r.Context())
	if err != nil {
		catalogErr(w, err)
		return
	}
	if cats == nil {
		cats = []string{}
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *handlers) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := h.deps.Catalog.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		catalogErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *handlers) createProduct(w http.ResponseWriter, r *http.Request) {
	var p catalog.Product
	if !decodeJSON(w, r, &p) {
		return
	}
	created, err := h.deps.Catalog.Create(r.Context(), p)
	if err != nil {
		catalogErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *handlers) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p catalog.Product
	if !decodeJSON(w, r, &p) {
		return
	}
	p.ID = r.PathValue("id")
	updated, err := h.deps.Catalog.Update(r.Context(), p)
	if err != nil {
		catalogErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *handlers) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := h.deps.Catalog.Delete(r.Context(), r.PathValue("id")); err != nil {
		catalogErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// uploadImage stores a product photo (multipart field "image", max 5MB)
// and sets the product ImageURL to the public /static/ path.
func (h *handlers) uploadImage(w http.ResponseWriter, r *http.Request) {
	if h.deps.Media == nil {
		writeErr(w, http.StatusServiceUnavailable, "media store not configured")
		return
	}
	id := r.PathValue("id")
	url, err := h.deps.Media.SaveImage(id, r)
	if err != nil {
		switch {
		case errors.Is(err, media.ErrTooBig):
			writeErr(w, http.StatusRequestEntityTooLarge, err.Error())
		case errors.Is(err, media.ErrBadType), errors.Is(err, media.ErrNoFile):
			writeErr(w, http.StatusBadRequest, err.Error())
		default:
			writeErr(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	p, err := h.deps.Catalog.Get(r.Context(), id)
	if err != nil {
		catalogErr(w, err)
		return
	}
	p.ImageURL = url
	updated, err := h.deps.Catalog.Update(r.Context(), p)
	if err != nil {
		catalogErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// serveStatic serves upload files without directory listing.
func (h *handlers) serveStatic(w http.ResponseWriter, r *http.Request) {
	if h.deps.Media == nil {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/static/")
	if name == "" || strings.Contains(name, "/") || strings.HasPrefix(name, ".") {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
	http.ServeFile(w, r, filepath.Join(h.deps.Media.Dir(), name))
}

func (h *handlers) getCart(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	c, err := h.deps.Cart.Get(r.Context(), sessionID)
	if err != nil {
		cartErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type addItemReq struct {
	SessionID string `json:"session_id"`
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

func (h *handlers) addCartItem(w http.ResponseWriter, r *http.Request) {
	var req addItemReq
	if !decodeJSON(w, r, &req) {
		return
	}
	c, err := h.deps.Cart.AddItem(r.Context(), req.SessionID, req.ProductID, req.Qty)
	if err != nil {
		cartErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *handlers) removeCartItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	c, err := h.deps.Cart.RemoveItem(r.Context(), sessionID, r.PathValue("productID"))
	if err != nil {
		cartErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type setQtyReq struct {
	SessionID string `json:"session_id"`
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

func (h *handlers) setCartItemQty(w http.ResponseWriter, r *http.Request) {
	var req setQtyReq
	if !decodeJSON(w, r, &req) {
		return
	}
	c, err := h.deps.Cart.SetQty(r.Context(), req.SessionID, req.ProductID, req.Qty)
	if err != nil {
		cartErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type checkoutReq struct {
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
	Fail      bool   `json:"fail,omitempty"`
}

func (h *handlers) checkout(w http.ResponseWriter, r *http.Request) {
	var req checkoutReq
	if !decodeJSON(w, r, &req) {
		return
	}
	o, err := h.deps.Order.Checkout(r.Context(), req.SessionID, req.Email, req.Fail)
	if err != nil {
		// Payment failure still returns the order for retry visibility.
		if errors.Is(err, order.ErrChargeFail) && o.ID != "" {
			writeJSON(w, http.StatusPaymentRequired, o)
			return
		}
		orderErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, o)
}

func (h *handlers) getOrder(w http.ResponseWriter, r *http.Request) {
	o, err := h.deps.Order.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		orderErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// listMyOrders returns public order history for one email (?email=).
func (h *handlers) listMyOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.deps.Order.ListByEmail(r.Context(), r.URL.Query().Get("email"))
	if err != nil {
		orderErr(w, err)
		return
	}
	if orders == nil {
		orders = []order.Order{}
	}
	writeJSON(w, http.StatusOK, orders)
}

func (h *handlers) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.deps.Order.List(r.Context())
	if err != nil {
		orderErr(w, err)
		return
	}
	if orders == nil {
		orders = []order.Order{}
	}
	writeJSON(w, http.StatusOK, orders)
}

type adminStats struct {
	Products     int   `json:"products"`
	Orders       int   `json:"orders"`
	PaidOrders   int   `json:"paidOrders"`
	RevenueMinor int64 `json:"revenueMinor"`
	LowStock     int   `json:"lowStock"`
}

// adminStats aggregates shop numbers for the admin dashboard.
func (h *handlers) adminStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := adminStats{}
	for page := 1; ; page++ {
		res, err := h.deps.Catalog.List(ctx, catalog.ListParams{Page: page, PerPage: 100})
		if err != nil {
			catalogErr(w, err)
			return
		}
		stats.Products = res.Total
		for _, p := range res.Items {
			if p.Stock <= 5 {
				stats.LowStock++
			}
		}
		if len(res.Items) < res.PerPage {
			break
		}
	}
	orders, err := h.deps.Order.List(ctx)
	if err != nil {
		orderErr(w, err)
		return
	}
	stats.Orders = len(orders)
	for _, o := range orders {
		switch o.Status {
		case order.StatusPaid, order.StatusShipped, order.StatusDone:
			stats.PaidOrders++
			stats.RevenueMinor += o.TotalMinor
		}
	}
	writeJSON(w, http.StatusOK, stats)
}

type statusReq struct {
	Status order.Status `json:"status"`
}

func (h *handlers) setOrderStatus(w http.ResponseWriter, r *http.Request) {
	var req statusReq
	if !decodeJSON(w, r, &req) {
		return
	}
	o, err := h.deps.Order.SetStatus(r.Context(), r.PathValue("id"), req.Status)
	if err != nil {
		orderErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if !decodeJSON(w, r, &req) {
		return
	}
	sess, err := h.deps.Auth.Login(req.Email, req.Password)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": sess.Token, "role": string(sess.Role)})
}

func (h *handlers) me(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 {
		token = token[7:]
	}
	sess, err := h.deps.Auth.Authenticate(token)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"email": sess.Email, "role": string(sess.Role)})
}

// readyz reports serving readiness; pings Postgres when configured.
func (h *handlers) readyz(w http.ResponseWriter, r *http.Request) {
	if h.deps.Ready != nil {
		if err := h.deps.Ready(r.Context()); err != nil {
			writeErr(w, http.StatusServiceUnavailable, "not ready")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ready": true})
}
