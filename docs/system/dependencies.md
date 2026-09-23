# Dependencies — Open Commerce

## Component Dependency Graph

```text
browser --> frontend (SvelteKit) --> http-api (Go) --> catalog/cart/order services --> store (memory | Postgres via pgxpool)
                                                       order service --> payment-stub
                                                       http-api --> auth-stub (middleware)
                                                       api startup --> db migrate + seed (DATABASE_URL set)
```

## Runtime Dependencies

| Dependent | Depends On | Type | Notes |
|-----------|------------|------|-------|
| `app/backend` | Go stdlib (`net/http`, `encoding/json`, `log/slog`, `embed`) | internal | — |
| `app/backend` | `github.com/jackc/pgx/v5` v5.11 | third-party (allowlisted, ADR-0003) | Pool + migrasi + PG stores; satu-satunya deps luar stdlib |
| `app/frontend` | SvelteKit 2, Svelte 4, Vite 5, TypeScript 5 | third-party | `npm install`; versi dikunci di `package.json` |
| `app/frontend` | Go API `/api/v1/*` | internal | Kontrak JSON di `components.md` + `lib/api.ts`; admin bawa Bearer token |
| services | `Store` interfaces (ctx + error) | internal | Memory impl default; PG impl bila `DATABASE_URL` terisi |

## External Dependencies

| Dependency | Provider | Purpose | Contract | Fallback |
|------------|----------|---------|----------|----------|
| None (v0.1) | — | — | Payment/auth adalah stub in-repo | Integrasi real wajib sediakan fallback mock + ADR |

## Infrastructure Dependencies

| Dependent | Infrastructure | Notes |
|-----------|----------------|-------|
| Go API | `PORT` env (default 8080), `DATABASE_URL` (kosong → memory; terisi → Postgres + migrasi + seed) | `config/app.yaml` + env override; compose di `infra/` |
| SvelteKit web | `PUBLIC_API_BASE_URL` (default `http://localhost:8080`) | Build-time public env; prod via proxy same-origin |
| Local dev | Docker Compose `db` (Postgres 16) — aktif dipakai sejak v0.2 | `docker compose -f infra/compose.yaml up db`, lalu `DATABASE_URL=... go run ./cmd/api` |
| Prod (planned) | GCP via Terraform: VPC, Cloud SQL 16, Artifact Registry, Cloud Run ×2 | Scaffold `infra/terraform/` belum di-apply; secrets via Secret Manager (ADR-0004) |
