# Project — Open Commerce

## Name

Open Commerce

## Purpose

Aplikasi e-commerce open source single-vendor (satu penjual) yang menyediakan katalog produk, keranjang, checkout, dan admin operasional. Ada untuk memberi basis yang dapat di-fork, dipelajari, dan dikembangkan komunitas tanpa ketergantungan proprietary.

## Users

| User / Persona | Need | Notes |
|----------------|------|-------|
| Pembeli (customer) | Browse katalog, cari produk, kelola keranjang, checkout | Pengguna utama storefront |
| Admin toko | Kelola produk, harga/stok, lihat dan proses order | Operator tunggal, bukan multi-seller |
| Kontributor open source | Memahami arsitektur cepat, menjalankan lokal < 10 menit, menambah fitur aman | Fork-friendly, MIT |

## Scope

- Katalog produk (CRUD admin, list/detail/search publik).
- Keranjang belanja (session-based, server-side di v0.1 fondasi).
- Order + checkout dengan payment stub (mock, tanpa uang real).
- Auth stub (session/token sederhana, bukan OAuth prod).
- Storefront SvelteKit (katalog, cart, checkout, sukses order).
- Admin minimal (produk + order list).
- API REST `GET/POST /api/v1/*` untuk katalog/cart/order.
- Observability minimal: structured stdout log + `GET /healthz`.

## Non-Goals

- Multi-vendor / marketplace (seller onboarding, komisi, split payout).
- Integrasi payment real (Midtrans/Xendit/Stripe) — hanya interface stub di v0.1; integrasi real butuh ADR baru.
- Mobile native apps, real-time chat, rekomendasi ML.
- Multi-currency / multi-warehouse / ERP sync.

Non-goals are binding. Work that falls under non-goals requires an explicit decision change, not silent scope expansion.

## Current Status

- Status: `active`
- Maturity: `prototype` (v0.1 fondasi + skeleton)
- Additional notes: Fondasi diinisialisasi 2026-09-23. Roadmap: v0.2 katalog lengkap + Postgres, v0.3 cart/order persisten, v0.4 payment real (ADR).

## Ownership

- Owners: Community / Maintainer (Unassigned — open for contributors).
- Decision authority: Maintainer via ADR + PR review.
