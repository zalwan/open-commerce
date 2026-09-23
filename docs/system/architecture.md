# Architecture — Open Commerce

## Architectural Style

- Style: `modular monolith` (backend Go) + `client-server` (SvelteKit frontend → Go REST API).
- Key patterns: handler → service → store (backend); file-based routing + load functions (frontend); interface-stub untuk payment/auth.
- Decision references: [0001-go-sveltekit-modular-monolith](../decisions/0001-go-sveltekit-modular-monolith.md), [0002-stub-payment-auth](../decisions/0002-stub-payment-auth.md).

## System Overview

Single-vendor toko online: SvelteKit me-render storefront/admin dan memanggil Go API via REST JSON. Go API menegakkan aturan domain (katalog, cart, order) dan persist via interface `Store` (in-memory di v0.1, Postgres di v0.2+). Payment dan auth adalah stub in-process yang bisa diganti implementasi real tanpa mengubah handler.

```text
browser --> SvelteKit web (:5173/:3000) --REST /api/v1--> Go API (:8080) --> Store (memory v0.1 / Postgres v0.2+)
                                                     Go API --> PaymentStub (mock)
                                                     Go API --> AuthStub (mock token)
```

## Components

Summarized here; detail in [components.md](./components.md).

| Component | Responsibility | Location |
|-----------|----------------|----------|
| Catalog service | Produk CRUD, search/list publik | `app/backend/internal/catalog/` |
| Cart service | Keranjang per-session, tambah/ubah/hapus | `app/backend/internal/cart/` |
| Order service | Checkout → order, status lifecycle | `app/backend/internal/order/` |
| Auth stub | Login mock → token, middleware | `app/backend/internal/auth/` |
| Payment stub | Charge/refund mock, selalu sukses/kontrol via flag | `app/backend/internal/payment/` |
| HTTP API layer | Routing `/api/v1/*`, JSON encode, request log, `/healthz` | `app/backend/internal/http/` + `cmd/api/` |
| Storefront web | Katalog, detail, cart, checkout, sukses | `app/frontend/src/routes/` (+ `lib/api.ts`, `lib/cart.ts`) |
| Admin web | Produk list/form, order list/detail | `app/frontend/src/routes/admin/` |

## Communication Mechanisms

| From → To | Mechanism | Notes |
|-----------|-----------|-------|
| Browser → SvelteKit | Request-response HTTP (SSR + CSR hydration) | Sync; SvelteKit load() untuk data awal |
| SvelteKit → Go API | REST JSON `GET/POST /api/v1/*` | Sync; base URL via `PUBLIC_API_BASE_URL`, kontrak di `docs/system/components.md` + tipe TS di `lib/api.ts` |
| HTTP handler → Service | In-process function calls (Go) | Sync; handler tidak akses store langsung |
| Service → Store | In-process via `Store` interfaces | Sync; in-memory v0.1, Postgres v0.2 |
| Order service → Payment stub | In-process interface `Charge()` | Sync; mock latency kecil, status deterministik |

## Data Flows

1. Browse & beli: `GET /api/v1/products` → render katalog → `POST /api/v1/cart/items` (session) → `POST /api/v1/orders/checkout` → order service panggil payment stub → order `paid|pending` → frontend halaman sukses.
2. Admin kelola: login stub → token → `POST /api/v1/admin/products` → catalog service validasi → store simpan → list ter-refresh. Order: `GET /api/v1/admin/orders` → ubah status (`processing → shipped → done`).

## Boundaries

- Module/component boundaries: `internal/http` tidak boleh SQL / akses map store langsung — hanya via service. Frontend tidak boleh akses DB — hanya via REST. Lihat `components.md`.
- Trust boundaries: validasi input di handler (400 on bad JSON); otorisasi admin di middleware auth-stub (401/403); aturan bisnis (stok > 0, total = sum items) di service.
- External boundaries: hanya REST publik `/api/v1/*` + `/healthz`. Tidak ada outbound ke vendor di v0.1.

## External Integrations

| Integration | Direction | Contract | Failure Handling |
|-------------|-----------|----------|------------------|
| None (v0.1) | — | — | Payment stub tidak pernah timeout nyata; error disimulasikan via `fail=true` untuk uji path gagal |

## Security Boundaries

- Auth stub: Bearer token opaque in-memory, expiry 24 jam; admin routes butuh `role=admin`. Bukan keamanan prod — diganti OAuth/JWT real di ADR lanjutan.
- Tidak ada secret di repo; konfigurasi via env (`PORT`, `PUBLIC_API_BASE_URL`, `DATABASE_URL` untuk v0.2).
- CORS: dev `*` untuk kemudahan lokal; prod wajib same-origin / allowlist (lihat constraints).

## Scalability Considerations

- Stateless API kecuali in-memory store v0.1 (tidak scale horizontal — diketahui, diganti Postgres).
- SvelteKit SSR stateless; session cart di server Go (cookie session id) — pindah ke Redis/DB bila scale.
- Bottleneck diketahui: single-process Go + map + mutex cukup untuk fondasi/demo.

## Failure Modes

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Go API down | Storefront gagal fetch → halaman error dengan retry | `/healthz` untuk probe; frontend tampilkan pesan + tombol retry |
| Payment stub `fail=true` | Checkout 402 + order `payment_failed` | Frontend tampilkan error; user bisa retry; tidak ada charge ganda (idempotency-key sederhana di v0.2) |
| Postgres belum ada (v0.1) | Tidak berdampak — memory store aktif | Migrasi v0.2 membawa `DATABASE_URL` + health check DB |
