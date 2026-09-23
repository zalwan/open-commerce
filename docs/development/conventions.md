# Engineering Conventions — Open Commerce

> Project-specific rules. Generic workflow: `.ai/workflow.md`, principles: `.ai/principles.md`.

## Code Style

- Go: `gofmt -l .` harus bersih; paket `internal/{catalog,cart,order,auth,payment,http}`; handler tidak akses store langsung (via service). Harga int64 minor units, JSON snake? tidak — pakai camelCase sesuai `lib/api.ts` (`priceMinor`, `session_id` hanya untuk body cart/checkout warisan — konsisten dengan tipe TS).
- SvelteKit: TypeScript `strict`, file-routing di `src/routes`, API client hanya via `src/lib/api.ts`, session/token via `src/lib/session.ts`. Tidak ada `any` tanpa komentar alasan.
- Format: Go `gofmt`, frontend `npx prettier --check` bila ditambah (opsional v0.1).

## Branching and Commits

- Branch: `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, `infra/<slug>`.
- Commit: `<type>: <ringkas>` — type: feat/fix/docs/refactor/infra/test/chore. Contoh: `feat: cart checkout stub`.
- PR wajib: `go test ./...` + `npm run build` hijau; sertakan impact areas (workflow §3).

## Testing

- Backend: `go test ./...` di `app/backend`. Service baru wajib unit test (validasi + happy/fail path). HTTP route baru wajib test status.
- Frontend: `npm run build` wajib lolos; `npm run check` bila disentuh routes/lib.
- E2E smoke: `sh tests/smoke.sh` (butuh API di :8080).
- Coverage target v0.1: domain services teruji; HTML routes manual.

## Configuration and Secrets

- Config berlapis: `config/app.yaml` (default) < env (`PORT`, `DATABASE_URL`, `PUBLIC_API_BASE_URL`).
- Dilarang commit `.env` berisi secret, `*.pem/*.key`, dump DB. Kredensial di repo hanya dummy yang didokumentasikan.
- Tambah env baru → update `config/app.yaml` + `setup.md` + `compose.yaml` bila relevan.

## Documentation

- Ubah arsitektur/komponen → update `docs/system/architecture.md` + `components.md`.
- Ganti/tambah teknologi → `technology.md` + `.ai/project.yaml`.
- Keputusan arsitektural (gaya, batas, integrasi, payment/auth real) → ADR baru `docs/decisions/NNNN-*.md`; jangan edit ADR lama selain typo/status.
- Operasional (deploy/monitor) → `docs/operations/` + `infra/` selaras.
