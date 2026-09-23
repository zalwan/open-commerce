# Technology — Open Commerce

> Authoritative inventory. List only what the project actually uses.

## Languages

| Language | Version | Purpose | Notes |
|----------|---------|---------|-------|
| Go | 1.25 | Backend API (catalog/cart/order/auth/payment-stub) | `app/backend/`, stdlib-first |
| TypeScript | 5.x | Frontend SvelteKit (storefront + admin) | `app/frontend/`, strict mode |
| JavaScript (Svelte) | Svelte 4 / SvelteKit 2 | UI components & routing | SSR + CSR hybrid |

## Frameworks and Libraries

| Name | Version | Purpose | Notes |
|------|---------|---------|-------|
| Go stdlib `net/http` | 1.25 | HTTP routing & handlers | Tanpa framework eksternal di v0.1 agar skeleton zero-deps |
| SvelteKit | 2.x | Frontend app framework, file-based routing | `@sveltejs/kit`, adapter-node untuk deploy |
| Vite | 5.x | Frontend dev server & build | Bawaan SvelteKit |

## Runtimes

| Runtime | Version | Purpose | Notes |
|---------|---------|---------|-------|
| Go toolchain | 1.25 | Build & run backend | `go build ./...`, `go test ./...` |
| Node.js | 22 LTS | Dev & build frontend, SSR runtime | `npm run dev/build/preview` |

## Data Stores

| Store | Type | Purpose | Notes |
|-------|------|---------|-------|
| Postgres | Relational | Prod store untuk produk, cart, order, user (doc intent) | v16; v0.1 skeleton pakai in-memory `Store` di balik interface agar bisa jalan tanpa DB |
| In-memory maps | Ephemeral (process-local) | Skeleton persistence v0.1 | Diganti Postgres tanpa ubah handler/service (lihat ADR-0001) |

## Infrastructure and Hosting

| Area | Choice | Notes |
|------|--------|-------|
| Hosting | Docker Compose lokal (api + web + db); image generik untuk prod | Lihat `infra/` |
| Networking | Browser → SvelteKit (`:5173` dev / `:3000` preview) → Go API (`:8080`) → Postgres (`:5432`) | REST JSON `/api/v1/*`, CORS dev open, prod same-origin via proxy |
| Storage | Postgres volume `pgdata`; upload produk out-of-scope v0.1 (URL eksternal saja) | — |

## CI/CD

| Stage | Tool / Mechanism | Notes |
|-------|------------------|-------|
| Build | `go build ./...` + `npm run build` (frontend) | GitHub Actions `.github/workflows/` |
| Test | `go test ./...` + `npm run check` (bila tersedia) | Wajib hijau sebelum merge |
| Deploy | Docker build (api, web) + compose pull | Manual di v0.1; pipeline deploy otomatis non-goal |

## Observability

| Area | Tool / Mechanism | Notes |
|------|------------------|-------|
| Logs | Structured stdout (Go `log/slog`, frontend console) | Request log: method path status latency |
| Metrics | None (v0.1) | `/healthz` untuk liveness saja |
| Traces | None | — |

## External Services

| Service | Purpose | Integration | Notes |
|---------|---------|-------------|-------|
| None (payment/auth stub in-repo) | Simulasi checkout & login tanpa vendor | In-process Go interfaces `PaymentProvider`, `Authenticator` | Integrasi real (Midtrans/Stripe/OAuth) butuh ADR baru + fallback mock |
