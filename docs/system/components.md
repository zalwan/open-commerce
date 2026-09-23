# Components — Open Commerce

## Catalog Service

- **Responsibility:** Aturan produk (nama wajib, harga > 0, stok >= 0); list/search publik, CRUD admin. Tidak mengurus cart/order.
- **Source location:** `app/backend/internal/catalog/`
- **Dependencies:** `Store` interface (in-memory v0.1).
- **Consumers:** HTTP handlers `GET /api/v1/products`, `GET /api/v1/products/{id}`, admin product endpoints.
- **Interfaces:** `Service.List(q string)`, `Service.Get(id)`, `Service.Create/Update/Delete` — JSON di `internal/http/`.
- **Data stores:** Product map (memory) → Postgres `products` di v0.2.
- **Operational notes:** Read-heavy; validasi harga/stok di service, bukan handler.

---

## Cart Service

- **Responsibility:** Keranjang per-session (`session_id` cookie): add/update/remove/clear, hitung subtotal. Tidak persist lintas restart di v0.1.
- **Source location:** `app/backend/internal/cart/`
- **Dependencies:** Catalog service (cek produk & harga snapshot), `Store`.
- **Consumers:** `POST /api/v1/cart/items`, `GET /api/v1/cart`, `DELETE /api/v1/cart/items/{productID}`.
- **Interfaces:** `Service.Get(sessionID)`, `Service.AddItem(sessionID, productID, qty)`, dsb.
- **Data stores:** Cart map (memory) → Postgres/Redis di v0.3.
- **Operational notes:** Qty <= stok saat add (best-effort, bukan reservasi stok).

---

## Order Service

- **Responsibility:** Checkout: snapshot cart → hitung total → panggil payment stub → buat order dengan status (`pending|paid|payment_failed|shipped|done|cancelled`). Transisi status tervalidasi.
- **Source location:** `app/backend/internal/order/`
- **Dependencies:** Cart service, Payment stub, `Store`.
- **Consumers:** `POST /api/v1/orders/checkout`, `GET /api/v1/orders/{id}`, admin order endpoints.
- **Interfaces:** `Service.Checkout(sessionID, req)`, `Service.Get(id)`, `Service.SetStatus(id, status)`.
- **Data stores:** Order map (memory) → Postgres `orders` + `order_items` di v0.2.
- **Operational notes:** Idempotency minimal via session clear setelah sukses; payment gagal → order tetap tercatat `payment_failed` untuk retry.

---

## Auth Stub

- **Responsibility:** Login mock (`admin@shop.test` / `admin123` → token admin; user lain → token customer). Validasi Bearer di middleware. Bukan auth prod.
- **Source location:** `app/backend/internal/auth/`
- **Dependencies:** None.
- **Consumers:** Middleware `RequireAuth` / `RequireAdmin`; `POST /api/v1/auth/login`, `GET /api/v1/auth/me`.
- **Interfaces:** `Login(email, password) (token, role)`, `Middleware(next)`.
- **Data stores:** Token map in-memory.
- **Operational notes:** Token opaque, expiry 24h; ganti JWT/OAuth via ADR tanpa ubah handler (interface tetap).

---

## Payment Stub

- **Responsibility:** Simulasi charge/refund deterministik. Sukses default; gagal bila `amount <= 0` atau `fail=true`.
- **Source location:** `app/backend/internal/payment/`
- **Dependencies:** None.
- **Consumers:** Order service saat checkout.
- **Interfaces:** `Provider.Charge(amountMinor int64, currency string, opts) (txID, error)`, `Provider.Refund(txID)`.
- **Data stores:** None (log transaksi in-memory untuk inspeksi).
- **Operational notes:** Latency tiruan < 50ms; integrasi real wajib implement interface ini (lihat ADR-0002).

---

## HTTP API Layer

- **Responsibility:** Routing, JSON encode/decode, validasi bentuk request, request log, CORS dev, `GET /healthz`.
- **Source location:** `app/backend/internal/http/`, `app/backend/cmd/api/main.go`
- **Dependencies:** Semua services (di-wire di `main.go`).
- **Consumers:** Frontend SvelteKit + API client eksternal.
- **Interfaces:** REST `GET/POST /api/v1/*` — kontrak:
  - `GET /healthz` → `{"ok":true}`
  - `GET /api/v1/products?q=` → `[{id,name,priceMinor,currency,stock}]`
  - `POST /api/v1/cart/items {"session_id","product_id","qty"}` → cart
  - `POST /api/v1/orders/checkout {"session_id","email","fail?"}` → order
  - `POST /api/v1/auth/login` → `{"token","role"}`
- **Data stores:** None langsung.
- **Operational notes:** Port via `PORT` (default 8080); JSON error envelope `{"error":"..."}`.

---

## Storefront Web

- **Responsibility:** Katalog, detail produk, cart drawer/page, checkout form, halaman sukses. Fetch ke Go API via `lib/api.ts`.
- **Source location:** `app/frontend/src/routes/`, `app/frontend/src/lib/`
- **Dependencies:** Go API (`PUBLIC_API_BASE_URL`).
- **Consumers:** Pembeli (browser).
- **Interfaces:** Routes `/`, `/products/[id]`, `/cart`, `/checkout`, `/orders/[id]`; lib `api.ts` (typed fetch), `cart.ts` (session_id di localStorage).
- **Data stores:** None langsung (via API).
- **Operational notes:** SSR dengan fallback CSR bila API down (tampilkan retry).

---

## Admin Web

- **Responsibility:** Login stub, produk list/form, order list/detail + ubah status.
- **Source location:** `app/frontend/src/routes/admin/`
- **Dependencies:** Go API + token di localStorage.
- **Consumers:** Admin toko (browser).
- **Interfaces:** Routes `/admin`, `/admin/products`, `/admin/orders`; header `Authorization: Bearer <token>`.
- **Data stores:** None langsung.
- **Operational notes:** Guard client-side + server-side 401/403 dari API.
