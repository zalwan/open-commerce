# Components — Open Commerce

## Catalog Service

- **Responsibility:** Aturan produk (nama wajib, harga > 0, stok >= 0); list/search publik, CRUD admin. Tidak mengurus cart/order.
- **Source location:** `app/backend/internal/catalog/` (`service.go`, `memory.go`, `store_pg.go`)
- **Dependencies:** `Store` interface (memory default; `PGStore` via pgxpool bila `DATABASE_URL` terisi).
- **Consumers:** HTTP handlers `GET /api/v1/products`, `GET /api/v1/products/{id}`, admin product endpoints.
- **Interfaces:** `Service.List(ctx, ListParams{q,category,sort,page,perPage}) → ListResult{items,total,page,perPage}`, `Get`, `Create/Update/Delete`, `Categories`, `DecrementStock/IncrementStock` (semua ber-`error`) — JSON di `internal/http/`.
- **Data stores:** Memory map → tabel Postgres `products` (+ `SeedIfEmpty` untuk DB segar; migrasi `0002_category.sql`). Decrement atomik (`UPDATE ... WHERE stock >= qty`) agar checkout konkuren tidak oversell.
- **Operational notes:** Read-heavy; validasi harga/stok di service, bukan handler.

---

## Cart Service

- **Responsibility:** Keranjang per-session (`session_id` cookie): add/update/remove/clear, hitung subtotal. Tidak persist lintas restart di v0.1.
- **Source location:** `app/backend/internal/cart/` (`service.go`, `store_pg.go`)
- **Dependencies:** Catalog service (cek produk & harga snapshot), `Store` (memory/`carts` JSONB).
- **Consumers:** `POST /api/v1/cart/items`, `GET /api/v1/cart`, `DELETE /api/v1/cart/items/{productID}`.
- **Interfaces:** `Service.Get/AddItem/SetQty/RemoveItem/Clear` (ctx + error); cart hilang → dianggap kosong; qty 0 = hapus.
- **Data stores:** Memory map → tabel `carts(session_id, items JSONB)`.
- **Operational notes:** Qty <= stok saat add (best-effort, bukan reservasi stok).

---

## Order Service

- **Responsibility:** Checkout: snapshot cart → reservasi stok (decrement atomik) → panggil payment stub → buat order dengan status (`pending|paid|payment_failed|shipped|done|cancelled`). Gagal bayar → restock kompensasi + order `payment_failed`. Cancel dari `paid` → restock.
- **Source location:** `app/backend/internal/order/` (`service.go`, `store_pg.go`)
- **Dependencies:** Cart service, Payment stub, `Store`.
- **Consumers:** `POST /api/v1/orders/checkout`, `GET /api/v1/orders/{id}`, admin order endpoints.
- **Interfaces:** `Service.Checkout(ctx, sessionID, req)`, `Get`, `List`, `ListByEmail` (riwayat publik per email), `SetStatus` (lifecycle tervalidasi).
- **Data stores:** Memory map → tabel `orders` (items JSONB, ID dari sequence `order_seq` → `o-<n>`).
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

## Media Store

- **Responsibility:** Simpan foto produk (multipart `image`, maks 5MB, jpeg/png/webp/gif; sniff content, bukan ekstensi) dengan nama aman dari ID produk; serve publik tanpa directory listing.
- **Source location:** `app/backend/internal/media/`
- **Dependencies:** Local disk (`DATA_DIR`, default `./data/uploads`).
- **Consumers:** `POST /api/v1/admin/products/{id}/image` → set `Product.ImageURL` ke `/static/<id>.<ext>`.
- **Interfaces:** `SaveImage(productID, req) (url, error)`; static `GET /static/<file>`.
- **Data stores:** Filesystem (volume `apidata` di compose).
- **Operational notes:** Ganti file se-ID menimpa versi lama; backup via volume, bukan DB.

---

## DB Bootstrap

- **Responsibility:** Buka pool pgx, apply migrasi embed, seed-on-empty. Hanya aktif bila `DATABASE_URL` terisi.
- **Source location:** `app/backend/internal/db/` (+ `migrations/0001_init.sql`)
- **Dependencies:** Postgres 16 reachable.
- **Consumers:** `cmd/api/main.go` saat startup; `GET /readyz` untuk ping.
- **Interfaces:** `db.Open(ctx, url)`, `db.Migrate(ctx, pool)`.
- **Data stores:** Skema `products`, `carts`, `orders`, sequence `order_seq`.
- **Operational notes:** DDL aditif saja; startup gagal fast bila DB unreachable.

---

## HTTP API Layer

- **Responsibility:** Routing, JSON encode/decode, validasi bentuk request, pemetaan error domain→4xx (error store→500 generik), request log, CORS dev, `GET /healthz` + `GET /readyz`.
- **Source location:** `app/backend/internal/http/`, `app/backend/cmd/api/main.go`
- **Dependencies:** Semua services (di-wire di `main.go`).
- **Consumers:** Frontend SvelteKit + API client eksternal.
- **Interfaces:** REST `GET/POST/PUT/DELETE /api/v1/*` — kontrak:
  - `GET /healthz` → `{"ok":true}` (tanpa auth, tanpa DB)
  - `GET /readyz` → `{"ready":true}` atau 503 bila DB unreachable
  - `GET /api/v1/products?q=&category=&sort=&page=&per_page=` → `{items,total,page,perPage}`
  - `GET /api/v1/categories` → `[..]`
  - `GET /api/v1/cart?session_id=` → cart; `POST /api/v1/cart/items` (tambah); `PUT /api/v1/cart/items` (set qty, 0 = hapus); `DELETE /api/v1/cart/items/{id}?session_id=`
  - `POST /api/v1/orders/checkout {"session_id","email","fail?"}` → order (402 + body order bila payment gagal)
  - `GET /api/v1/orders?email=` → riwayat publik; `GET /api/v1/orders/{id}` → detail
  - `GET /static/<file>` → foto produk
  - `POST /api/v1/auth/login` → `{"token","role"}`
  - Admin (Bearer admin): CRUD produk + upload foto + `GET /api/v1/admin/orders` + `POST /api/v1/admin/orders/{id}/status` + `GET /api/v1/admin/stats`
- **Data stores:** None langsung.
- **Operational notes:** Port via `PORT` (default 8080); JSON error envelope `{"error":"..."}`.

---

## Storefront Web

- **Responsibility:** Katalog (filter kategori, pencarian, sort, pagination), detail produk, cart drawer/page dengan stepper qty, checkout + ringkasan, halaman sukses, lacak order per email. Fetch ke Go API via `lib/api.ts`.
- **Source location:** `app/frontend/src/routes/`, `app/frontend/src/lib/`
- **Dependencies:** Go API (`PUBLIC_API_BASE_URL`).
- **Consumers:** Pembeli (browser).
- **Interfaces:** Routes `/`, `/products/[id]`, `/cart`, `/checkout`, `/orders/[id]`, `/orders/track`; lib `api.ts` (typed fetch), `session.ts` (session_id + token di localStorage).
- **Data stores:** None langsung (via API).
- **Operational notes:** SSR dengan fallback CSR bila API down (tampilkan retry).

---

## Admin Web

- **Responsibility:** Dashboard stats (produk, order, revenue, low-stock), login stub, produk list + form tambah/edit + hapus + upload foto, order list/detail + advance status.
- **Source location:** `app/frontend/src/routes/admin/` (+ `lib/api.ts` fungsi admin + `uploadImage` FormData dengan header Bearer)
- **Dependencies:** Go API + token di localStorage.
- **Consumers:** Admin toko (browser).
- **Interfaces:** Routes `/admin`, `/admin/products`, `/admin/orders`, `/admin/orders/[id]`; guard: tanpa token → pesan login (server tetap 401/403).
- **Data stores:** None langsung.
- **Operational notes:** Guard client-side + server-side 401/403 dari API.
