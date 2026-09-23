# Constraints — Open Commerce

## Technical Constraints

- Go >= 1.25, Node >= 22 LTS. Versi minimum mengikat; upgrade butuh uji `go test` + `npm run build`.
- API versioning: publik hanya `/api/v1/*`. Breaking change butuh versi baru + ADR, tidak boleh ubah v1 in-place.
- Backend dependency third-party di-allowlist: hanya `github.com/jackc/pgx/v5` (ADR-0003). Tambah dependency lain butuh ADR + update `technology.md` + `project.yaml`.
- Migrasi DB hanya aditif (`CREATE TABLE/INDEX IF NOT EXISTS`, sequence). Alter destruktif (DROP/RENAME kolom) butuh ADR + rencana migrasi data.
- Frontend: TypeScript strict, tidak ada `any` tanpa alasan; SvelteKit file-routing tetap dipakai (jangan custom router).
- Data: harga dalam minor units (`priceMinor` int64, mis. rupiah sen) — tidak ada float untuk uang.

## Business Constraints

- Lisensi MIT — semua kontribusi harus kompatibel MIT; tidak boleh dependensi berlisensi copyleft kuat (GPL/AGPL) tanpa ADR + persetujuan maintainer.
- Single-vendor only di v0.1. Fitur multi-seller/komisi dilarang masuk tanpa ADR yang mencabut non-goal.
- Payment stub bukan alat transaksi real — dilarang memproses uang asli via stub.

## Operational Constraints

- `GET /healthz` (tanpa DB) dan `GET /readyz` (ping DB bila dikonfigurasi) wajib ada dan tanpa auth; dipakai probe Docker/K8s.
- Konfigurasi via env + `config/app.yaml`; dilarang hardcode port/URL/credential di kode.
- Upload file produk out-of-scope v0.1 — hanya URL gambar eksternal.

## Security Constraints

- Tidak ada secret/credential real di repo, issue, atau log. Token stub hanya untuk dev (`admin123` didokumentasikan sebagai kredensial dummy).
- Auth stub bukan keamanan prod — dilarang mengklaimnya aman untuk prod; endpoint admin wajib 401/403 bila token hilang/salah.
- CORS dev boleh `*`; prod wajib allowlist/same-origin. Validasi input di handler (bad JSON → 400).
- Secret handling: baca dari env; contoh nilai hanya dummy di `config/` dan `setup.md`.

## Compatibility Constraints

- REST JSON UTF-8; currency default `IDR`; format error `{"error":"..."}` konsisten.
- Frontend harus jalan dengan `npm run dev` tanpa DB; backend harus jalan dengan `go run ./...` tanpa DB (memory store).
- Browser support: evergreen 2 versi terakhir (SvelteKit default).

## Performance Constraints

- p95 checkout lokal < 500ms (stub, tanpa DB) — diukur manual di v0.1, benchmark formal di v0.3.
- Halaman katalog SSR < 2s di dev lokal (tanpa optimasi prod).

## Forbidden Changes

- Dilarang handler Go mengakses store/DB langsung — wajib via service.
- Dilarang frontend mengakses DB langsung — wajib via REST API.
- Dilarang menambah integrasi payment/auth real, vector DB, atau platform deploy baru tanpa ADR.
- Dilarang commit secret, `.env` berisi credential real, atau database dump.
- Dilarang mengubah `docs/decisions/*` yang sudah ada selain typo/status-line; perubahan butuh ADR baru.
- Dilarang memecah backend menjadi microservices tanpa ADR (tetap modular monolith di v0.1).
