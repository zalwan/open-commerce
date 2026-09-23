package httpapi

import (
	"net/http"

	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/order"
)

func (h *handlers) listProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	writeJSON(w, http.StatusOK, h.deps.Catalog.List(q))
}

func (h *handlers) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := h.deps.Catalog.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *handlers) createProduct(w http.ResponseWriter, r *http.Request) {
	var p catalog.Product
	if !decodeJSON(w, r, &p) {
		return
	}
	created, err := h.deps.Catalog.Create(p)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
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
	updated, err := h.deps.Catalog.Update(p)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *handlers) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := h.deps.Catalog.Delete(r.PathValue("id")); err != nil {
		writeErr(w, http.StatusNotFound, "product not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) getCart(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeErr(w, http.StatusBadRequest, "session_id required")
		return
	}
	writeJSON(w, http.StatusOK, h.deps.Cart.Get(sessionID))
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
	c, err := h.deps.Cart.AddItem(req.SessionID, req.ProductID, req.Qty)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *handlers) removeCartItem(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeErr(w, http.StatusBadRequest, "session_id required")
		return
	}
	writeJSON(w, http.StatusOK, h.deps.Cart.RemoveItem(sessionID, r.PathValue("productID")))
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
	o, err := h.deps.Order.Checkout(req.SessionID, req.Email, req.Fail)
	if err != nil {
		// Payment failure still returns the order for retry visibility.
		if o.ID != "" {
			writeJSON(w, http.StatusPaymentRequired, o)
			return
		}
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, o)
}

func (h *handlers) getOrder(w http.ResponseWriter, r *http.Request) {
	o, err := h.deps.Order.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "order not found")
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func (h *handlers) listOrders(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.deps.Order.List())
}

type statusReq struct {
	Status order.Status `json:"status"`
}

func (h *handlers) setOrderStatus(w http.ResponseWriter, r *http.Request) {
	var req statusReq
	if !decodeJSON(w, r, &req) {
		return
	}
	o, err := h.deps.Order.SetStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
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
