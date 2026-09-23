# 0003: Postgres sebagai store persisten via pgx, migrasi SQL ter-embed

- Status: Accepted
- Date: 2026-09-23
- Deciders: Maintainer (lanjutan v0.1 → v0.2: persistensi + admin penuh)

## Context

v0.1 menyimpan segalanya di map in-memory: hilang saat restart, tidak bisa di-scale, admin CRUD tanpa persistensi berarti. `config/app.yaml` dan compose sudah mendeklarasikan Postgres 16 sebagai intent (standby). Fase v0.2 membutuhkan persistensi nyata untuk katalog/cart/order tanpa mengubah kontrak REST `/api/v1/*` dan tanpa memecah modular monolith.

## Decision

- Postgres 16 sebagai store persisten; driver `github.com/jackc/pgx/v5` (pool native `pgxpool`, bukan `database/sql` shim) — satu-satunya dependency third-party backend.
- Interface `Store` di tiap domain (`catalog`, `cart`, `order`) diubah ke bentuk ber-`error` + `context.Context`; implementasi `Memory*` (dev/test tanpa DB) dan `PG*` (prod/lokal dengan DB) hidup berdampingan.
- Migrasi SQL ter-embed (`app/backend/internal/db/migrations/*.sql`, `embed.FS`), dijalankan saat startup bila `DATABASE_URL` terisi; hanya DDL aditif (`CREATE TABLE IF NOT EXISTS`, sequence). Tanpa tool migrasi eksternal di v0.2.
- `DATABASE_URL` kosong → memory + seed (aliran dev v0.1 tidak berubah). Terisi → Postgres; `GET /readyz` memeriksa konektivitas DB.
- Item cart/order disimpan sebagai `JSONB` snapshot (bukan tabel ternormalisasi penuh) — tradeoff v0.2 agar checkout tetap satu round-trip; normalisasi (`order_items`) ditunda ke v0.3 bila query analitik dibutuhkan.
- ID order dari sequence Postgres (`order_seq`) dengan prefix `o-`, menggantikan counter atomik in-memory.

## Alternatives Considered

- `database/sql` + `lib/pq` — ditolak: `lib/pq` dalam maintenance-mode; `pgx` adalah driver modern yang direkomendasikan komunitas dan mendukung pool + typed arrays tanpa shim.
- ORM (GORM/ent) — ditolak: overkill untuk 3 tabel; SQL eksplisit lebih mudah di-review kontributor dan selaras dengan prinsip smallest-correct-change.
- Tool migrasi eksternal (golang-migrate/atlas) — ditolak untuk v0.2: satu file DDL aditif cukup; dievaluasi ulang saat skema butuh alter destruktif.
- Normalisasi penuh `order_items` sejak awal — ditolak untuk v0.2: menambah join tanpa kebutuhan query; JSONB snapshot sudah memenuhi kontrak API yang ada.

## Consequences

- Plus: restart-safe, siap compose penuh (`api` + `db`), kontrak API tidak berubah, frontend/admin tidak perlu penyesuaian protokol.
- Minus/risiko: interface `Store` berubah (refactor service + test v0.1 dalam commit yang sama); JSONB mengorbankan query item-level di SQL — diterima, tercatat untuk v0.3.
- Perlu: update `technology.md` (pgx aktif), `project.yaml`, `architecture.md` (paragraf store), `components.md` (impl PG), `constraints.md` (cabut zero-deps → allowlist `pgx`), `setup.md` (aliran dengan DB), CI job Postgres + integration test bertag.

## Related

- ADR-0001 (modular monolith, handler→service→store), ADR-0002 (stub tetap in-process).
- `docs/system/technology.md`, `docs/system/architecture.md`, `infra/compose.yaml`, `config/app.yaml`.
