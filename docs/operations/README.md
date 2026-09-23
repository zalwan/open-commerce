# Operations — Open Commerce (v0.1 fondasi)

## Deployment

- Lokal: backend `go run ./cmd/api` (di `app/backend`, `:8080`) + frontend `npm run dev` (di `app/frontend`, `:5173`).
- Compose: `docker compose -f infra/compose.yaml up --build` (api :8080, web :3000, db Postgres standby).
- Konfigurasi: `PORT`, `DATABASE_URL` (opsional v0.1), `PUBLIC_API_BASE_URL`. Sumber: [`config/app.yaml`](../../config/app.yaml).
- Rollback: v0.1 tanpa migrasi DB — rollback = checkout commit sebelumnya + rebuild. Mulai v0.2 (Postgres) butuh `migrate down` terdokumentasi.

## Monitoring

- Health: `GET /healthz` → `{"ok":true}` (tanpa auth). Compose `api.healthcheck` + probe K8s memakai path ini.
- Logs: backend JSON via `log/slog` (method, path, latency_ms); frontend console. Cari `level=ERROR` untuk 5xx.
- Metrik/Trace: none di v0.1 (diterima untuk prototype).

## Troubleshooting

| Symptom | Likely Cause | Diagnosis | Remediation |
|---------|--------------|-----------|-------------|
| Katalog kosong / fetch gagal di :5173 | API :8080 mati atau `PUBLIC_API_BASE_URL` salah | `curl localhost:8080/api/v1/products` | Nyalakan API; samakan env |
| Checkout 402 `payment failed` | `fail=true` terkirim atau amount invalid | Cek body request; log API | Retry tanpa `fail`; pastikan cart tidak kosong |
| Admin 401/403 | Token hilang / bukan admin | `POST /api/v1/auth/login` dengan dummy admin | Login ulang, kirim `Authorization: Bearer <token>` |
| Port bentrok | Proses lama masih jalan | `lsof -i :8080` | Kill atau `PORT=8081` + sesuaikan frontend env |
