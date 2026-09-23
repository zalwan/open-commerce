# 0001: Go backend + SvelteKit frontend sebagai modular monolith

- Status: Accepted
- Date: 2026-09-23
- Deciders: Maintainer (user request: backend Go, frontend SvelteKit, single-vendor)

## Context

Butuh fondasi e-commerce open source single-vendor yang mudah di-fork: katalog, cart, checkout, admin. Repo masih template AI-Native kosong. User memilih backend Go dan frontend SvelteKit, scope fondasi + skeleton, payment/auth stub.

## Decision

- Backend Go 1.25 modular monolith (stdlib `net/http`), struktur `cmd/api` + `internal/{catalog,cart,order,auth,payment,http,store}`.
- Frontend SvelteKit 2 + TypeScript (`app/frontend`), panggil backend via REST `/api/v1/*`.
- Persistensi di balik interface `Store`: in-memory di v0.1, Postgres di v0.2+ tanpa ubah handler/service.
- API versioned `/api/v1`, harga dalam minor units, error envelope `{"error":"..."}`.

## Alternatives Considered

- Python FastAPI + Next.js — ditolak: tidak sesuai permintaan eksplisit user (Go + SvelteKit).
- Microservices per-domain — ditolak: overkill untuk single-vendor fondasi; operasional berat untuk kontributor.
- Monolit penuh tanpa frontend terpisah (Go render HTML) — ditolak: SvelteKit memberi UX storefront modern + pemisahan client/server yang diinginkan.

## Consequences

- Plus: zero-deps backend mudah di-review; frontend modern; batas modul jelas; migrasi ke Postgres tanpa rewrite.
- Minus/risiko: in-memory tidak persist & tidak scale — diterima untuk v0.1, dimigrasi di v0.2; dua runtime (Go + Node) untuk dev lokal.
- Perlu: `docs/system/*` diisi, `project.yaml` metadata, compose lokal, CI `go test` + `npm build`.

## Related

- `docs/system/architecture.md`, `docs/system/technology.md`, `docs/system/components.md`
- ADR berikutnya: 0002 stub payment/auth.
