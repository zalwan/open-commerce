# Development Setup — Open Commerce

## Prerequisites

- Go >= 1.25 (`go version`)
- Node.js >= 22 LTS (`node --version`) + npm
- (Opsional, untuk v0.2+) Docker + Docker Compose untuk Postgres 16

## Setup Steps

1. Clone lalu jalankan backend tanpa DB (mode memory + seed):
   ```sh
   go run ./cmd/api
   # di direktori app/backend; dengar di :8080, cek curl localhost:8080/healthz
   ```
2. Atau dengan Postgres (mode persisten, disarankan untuk dev fitur):
   ```sh
   docker compose -f infra/compose.yaml up db
   DATABASE_URL='postgres://opencommerce:opencommerce@localhost:5432/opencommerce?sslmode=disable' go run ./cmd/api
   # migrasi auto-apply, seed-on-empty, cek curl localhost:8080/readyz
   ```
3. Di terminal lain, jalankan frontend:
   ```sh
   npm install
   npm run dev -- --port 5173
   # di direktori app/frontend; buka http://localhost:5173
   ```
4. Konfirmasi: katalog tampil, tambah ke cart, checkout dengan email bebas → halaman order sukses. Coba admin: login `admin@shop.test / admin123` di `/admin`, tambah produk, ubah status order.

## Configuration

- Backend: `PORT` (default `8080`), `DATABASE_URL` (kosong → memory store), `DATA_DIR` (default `./data/uploads` untuk foto produk). Defaults di [`config/app.yaml`](../../config/app.yaml).
- Frontend: `PUBLIC_API_BASE_URL` (default `http://localhost:8080`).
- Kredensial dummy dev-only: `admin@shop.test / admin123`, `customer@shop.test / customer123`. Bukan untuk prod.

## Running

- Backend tests: `go test ./...` di `app/backend`; integration (butuh Postgres): `DATABASE_URL=... go test -tags integration ./...`.
- Frontend build: `npm run build` di `app/frontend`; type-check: `npm run check`.
- Smoke E2E: `sh tests/smoke.sh` (memory) atau `DATABASE_URL=... BASE=... sh tests/smoke.sh` (Postgres).

## Troubleshooting Setup

| Problem | Fix |
|---------|-----|
| `localhost:5173` kosong / fetch gagal | Pastikan API di `:8080` jalan; cek `curl localhost:8080/api/v1/products` |
| `npm run check` gagal `svelte-kit sync` | Jalankan `npm install` ulang; pastikan Node 22 |
| Port 8080 dipakai | `PORT=8081 go run ./cmd/api` + set `PUBLIC_API_BASE_URL=http://localhost:8081` saat `npm run dev` |
| Integration test skip (`DATABASE_URL not set`) | Ekspor `DATABASE_URL` ke Postgres lokal/compose dulu |
| `readyz` 503 | DB unreachable — cek compose `db` jalan dan kredensial `DATABASE_URL` |
