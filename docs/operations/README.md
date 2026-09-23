# Operations — Open Commerce (v0.2 Postgres + admin penuh)

## Deployment

- Lokal memory: backend `go run ./cmd/api` (di `app/backend`, `:8080`) + frontend `npm run dev` (di `app/frontend`, `:5173`).
- Lokal Postgres: `docker compose -f infra/compose.yaml up db`, lalu `DATABASE_URL='postgres://opencommerce:opencommerce@localhost:5432/opencommerce?sslmode=disable' go run ./cmd/api`.
- Compose penuh: `docker compose -f infra/compose.yaml up --build` (api :8080, web :3000, db Postgres; api auto-migrasi + seed-on-empty).
- Konfigurasi: `PORT`, `DATABASE_URL` (kosong → memory), `PUBLIC_API_BASE_URL`. Sumber: [`config/app.yaml`](../../config/app.yaml).
- Rollback: DDL v0.2 aditif sehingga rollback = checkout commit sebelumnya + rebuild tanpa migrasi turun. Data seed tidak dihapus otomatis.

## Monitoring

- Health: `GET /healthz` → `{"ok":true}` (tanpa auth/DB). Readiness: `GET /readyz` → `{"ready":true}` atau 503 bila Postgres unreachable. Compose `db.healthcheck` + probe K8s memakai keduanya.
- Logs: backend JSON via `log/slog` (method, path, latency_ms, plus `store=memory|postgres` saat start); frontend console. Cari `level=ERROR` untuk 5xx.
- Metrik/Trace: none di v0.2 (diterima untuk prototype).

## Troubleshooting

| Symptom | Likely Cause | Diagnosis | Remediation |
|---------|--------------|-----------|-------------|
| Katalog kosong / fetch gagal di :5173 | API :8080 mati atau `PUBLIC_API_BASE_URL` salah | `curl localhost:8080/api/v1/products` | Nyalakan API; samakan env |
| Checkout 402 `payment failed` | `fail=true` terkirim atau amount invalid | Cek body request; log API | Retry tanpa `fail`; pastikan cart tidak kosong |
| Admin 401/403 | Token hilang / bukan admin | `POST /api/v1/auth/login` dengan dummy admin | Login ulang, kirim `Authorization: Bearer <token>` |
| Port bentrok | Proses lama masih jalan | `lsof -i :8080` | Kill atau `PORT=8081` + sesuaikan frontend env |
| `readyz` 503 | Postgres unreachable | Cek `DATABASE_URL`, `docker compose … up db` | Darurat: kosongkan `DATABASE_URL` untuk fallback memory (data PG tidak hilang) |
