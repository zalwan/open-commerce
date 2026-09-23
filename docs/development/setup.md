# Development Setup — Open Commerce

## Prerequisites

- Go >= 1.25 (`go version`)
- Node.js >= 22 LTS (`node --version`) + npm
- (Opsional, untuk v0.2+) Docker + Docker Compose untuk Postgres 16

## Setup Steps

1. Clone lalu jalankan backend (tanpa DB di v0.1):
   ```sh
   go run ./cmd/api
   # di direktori app/backend; dengar di :8080, cek curl localhost:8080/healthz
   ```
2. Di terminal lain, jalankan frontend:
   ```sh
   npm install
   npm run dev -- --port 5173
   # di direktori app/frontend; buka http://localhost:5173
   ```
3. Konfirmasi: katalog tampil (3 produk seed), tambah ke cart, checkout dengan email bebas → halaman order sukses.
4. (Opsional) Nyalakan Postgres standby untuk v0.2:
   ```sh
   docker compose -f infra/compose.yaml up db
   ```

## Configuration

- Backend: `PORT` (default `8080`), `DATABASE_URL` (kosong di v0.1 → memory store). Defaults di [`config/app.yaml`](../../config/app.yaml).
- Frontend: `PUBLIC_API_BASE_URL` (default `http://localhost:8080`).
- Kredensial dummy dev-only: `admin@shop.test / admin123`, `customer@shop.test / customer123`. Bukan untuk prod.

## Running

- Backend tests: `go test ./...` di `app/backend`.
- Frontend build: `npm run build` di `app/frontend`; type-check: `npm run check`.
- Smoke E2E (butuh API jalan): `sh tests/smoke.sh`.

## Troubleshooting Setup

| Problem | Fix |
|---------|-----|
| `localhost:5173` kosong / fetch gagal | Pastikan API di `:8080` jalan; cek `curl localhost:8080/api/v1/products` |
| `npm run check` gagal `svelte-kit sync` | Jalankan `npm install` ulang; pastikan Node 22 |
| Port 8080 dipakai | `PORT=8081 go run ./cmd/api` + set `PUBLIC_API_BASE_URL=http://localhost:8081` saat `npm run dev` |
