# 0002: Payment dan auth sebagai stub in-repo

- Status: Accepted
- Date: 2026-09-23
- Deciders: Maintainer (user request: stub/mock saja)

## Context

Checkout butuh payment dan endpoint admin butuh auth, tapi integrasi real (Midtrans/Stripe/OAuth) out-of-scope fondasi dan berisiko (secret, compliance, uang asli).

## Decision

- `PaymentProvider` interface (`Charge/Refund`) dengan impl `Stub` deterministik: sukses default, gagal bila amount invalid atau `fail=true`.
- `Authenticator` stub: login mock (`admin@shop.test`/`admin123` → admin), token opaque in-memory, middleware `RequireAuth/RequireAdmin`.
- Integrasi real wajib implement interface yang sama + ADR baru + fallback mock; dilarang proses uang asli via stub.

## Alternatives Considered

- Langsung integrasi Midtrans/Stripe sandbox — ditolak: butuh secret, akun, webhook; menghambat quickstart < 10 menit.
- Tanpa payment sama sekali (order selalu sukses) — ditolak: path gagal tidak teruji; stub memberi kontrol deterministik untuk test.
- JWT/OAuth penuh di v0.1 — ditolak: kompleksitas tanpa manfaat fondasi; interface stub memungkinkan swap nanti.

## Consequences

- Plus: checkout end-to-end jalan lokal tanpa vendor; test deterministik; batas keamanan jelas (stub ≠ prod).
- Risiko: kontributor mengira stub aman untuk prod — dimitigasi via `constraints.md` + komentar kode + docs.
- Perlu: dokumentasikan kredensial dummy hanya untuk dev; tambah ADR saat payment/auth real dipilih.

## Related

- `docs/system/architecture.md` (failure modes), `docs/system/constraints.md` (security), ADR-0001.
