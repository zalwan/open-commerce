# Dependencies — Open Commerce

## Component Dependency Graph

```text
browser --> frontend (SvelteKit) --> http-api (Go) --> catalog/cart/order services --> store (memory v0.1 / Postgres v0.2+)
                                                       order service --> payment-stub
                                                       http-api --> auth-stub (middleware)
```

## Runtime Dependencies

| Dependent | Depends On | Type | Notes |
|-----------|------------|------|-------|
| `app/backend` | Go stdlib only (`net/http`, `encoding/json`, `log/slog`) | internal | Zero third-party deps di v0.1; tambah chi/pgx hanya via ADR |
| `app/frontend` | SvelteKit 2, Svelte 4, Vite 5, TypeScript 5 | third-party | `npm install`; versi dikunci di `package.json` |
| `app/frontend` | Go API `/api/v1/*` | internal | Kontrak JSON di `components.md` + `lib/api.ts`; base URL via env |
| services | `Store` interfaces | internal | In-memory impl v0.1; Postgres impl v0.2 tanpa ubah handler |

## External Dependencies

| Dependency | Provider | Purpose | Contract | Fallback |
|------------|----------|---------|----------|----------|
| None (v0.1) | — | — | Payment/auth adalah stub in-repo | Integrasi real wajib sediakan fallback mock + ADR |

## Infrastructure Dependencies

| Dependent | Infrastructure | Notes |
|-----------|----------------|-------|
| Go API | `PORT` env (default 8080), Postgres `DATABASE_URL` (v0.2+, opsional di v0.1) | `config/app.yaml` + env override; compose di `infra/` |
| SvelteKit web | `PUBLIC_API_BASE_URL` (default `http://localhost:8080`) | Build-time public env; prod via proxy same-origin |
| Local dev | Docker Compose `db` (Postgres 16) — standby untuk v0.2 | v0.1 jalan tanpa DB; `docker compose up db` opsional |
