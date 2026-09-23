# Open Commerce

An open source single-vendor e-commerce app (one seller): product catalog, cart, checkout, and operations admin. **Go** backend, **SvelteKit** frontend, **Postgres** database. MIT-licensed, fork-friendly.

## Features (v0.2)

- **Storefront:** catalog + search, product detail, per-session cart, checkout, order success page.
- **Stock safety:** atomic reservation at checkout (409 when insufficient), restock on payment failure or cancel-from-paid.
- **Admin:** login, product CRUD (create/edit/delete), order list, advance status (`paid → shipped → done`).
- **Versioned REST API** `/api/v1/*` — contract in [`docs/system/components.md`](docs/system/components.md).
- **In-repo payment & auth stubs** (deterministic mocks, no real money) — replaced by real integrations via ADR.
- **Two run modes:** memory (no DB, auto-seed) or Postgres (auto-applied migrations + seed-on-empty).
- **Health:** `GET /healthz` (liveness), `GET /readyz` (readiness, pings DB when configured).

## Tech Stack

| Layer | Choice |
|-------|--------|
| Backend | Go 1.25, stdlib `net/http`, `pgx/v5` (only external dep) |
| Frontend | SvelteKit 2, Svelte 4, strict TypeScript, Vite 5 |
| Database | Postgres 16 (active when `DATABASE_URL` is set), in-memory fallback |
| Local infra | Docker Compose (`api`, `web`, `db`) |
| CI | GitHub Actions: `go vet` + unit + integration + `npm run check` + build |

Architecture: modular monolith (handler → service → store) + client-server. Decisions recorded in [`docs/decisions/`](docs/decisions/) (ADR-0001 through 0003).

## Quickstart

Prerequisites: Go ≥ 1.25, Node ≥ 22. Optional: Docker for Postgres.

**Without a DB (fastest):**

```sh
# terminal 1 — API on :8080
cd app/backend && go run ./cmd/api

# terminal 2 — web on :5173
cd app/frontend && npm install && npm run dev
```

**With Postgres (persistent):**

```sh
docker compose -f infra/compose.yaml up db
cd app/backend && DATABASE_URL='postgres://opencommerce:opencommerce@localhost:5432/opencommerce?sslmode=disable' go run ./cmd/api
```

Open http://localhost:5173. Try: catalog → add to cart → checkout (any email) → success page. Admin at http://localhost:5173/admin with dev-only dummy credentials **`admin@shop.test` / `admin123`** (not for production).

Quick verification (requires the API running): `sh tests/smoke.sh`.

## Repo Structure

```text
app/backend/    Go API (cmd/api, internal/{catalog,cart,order,auth,payment,http,db})
app/frontend/   SvelteKit (catalog/cart/checkout/orders/admin routes, lib/api.ts)
config/         app.yaml — defaults (env overrides: PORT, DATABASE_URL, PUBLIC_API_BASE_URL)
infra/          compose.yaml, Dockerfile.api/web
tests/          smoke.sh — E2E (health, checkout, admin CRUD, 401 guard)
docs/system/    project, technology, architecture, components, dependencies, constraints
docs/decisions/ ADRs 0001–0003 (stack, stubs, Postgres)
```

## Configuration

| Env | Default | Notes |
|-----|---------|-------|
| `PORT` | `8080` | API port |
| `DATABASE_URL` | _(empty → memory)_ | Example: `postgres://opencommerce:opencommerce@localhost:5432/opencommerce?sslmode=disable` |
| `PUBLIC_API_BASE_URL` | `http://localhost:8080` | API base URL for the frontend |

Never commit secrets — only documented dummy credentials.

## Testing

```sh
cd app/backend && go vet ./... && go test ./... && go build ./...
# integration (requires Postgres):
DATABASE_URL='postgres://...' go test -tags integration ./...

cd app/frontend && npm run check && npm run build
```

## Roadmap

- v0.1 — foundation + skeleton ✅
- v0.2 — Postgres + full admin ✅ (current)
- v0.3 — candidates: real payment (sandbox), real auth (JWT/OAuth), `order_items` normalization, pagination/search, public deploy. Real payment/auth require an ADR.

Binding non-goals: multi-vendor/marketplace, native mobile, ERP sync (see [`docs/system/project.md`](docs/system/project.md)).

## Contributing

Read [`AGENTS.md`](AGENTS.md) (workflow), [`docs/development/setup.md`](docs/development/setup.md), and [`docs/development/conventions.md`](docs/development/conventions.md). Branch `feat/<slug>`; commits `<type>: <summary>`; PRs must be green (backend + frontend). Architectural changes need a new ADR — never edit old ADRs.

## License

MIT — see [`LICENSE`](LICENSE).
