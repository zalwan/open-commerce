# Backend — Open Commerce API (Go)

Modular monolith, stdlib `net/http`, zero third-party deps (v0.1).

## Run

```sh
go run ./cmd/api              # :8080
PORT=8081 go run ./cmd/api
curl localhost:8080/healthz
```

## Test / Build

```sh
go vet ./...
go test ./...
go build ./...
```

## Layout

- `cmd/api/main.go` — wiring + listen
- `internal/catalog` — product rules + memory store (seed 3 produk)
- `internal/cart` — session carts
- `internal/order` — checkout + lifecycle
- `internal/auth` — stub login (`admin@shop.test/admin123`)
- `internal/payment` — stub charge/refund (no real money)
- `internal/http` — routes `/api/v1/*`, CORS dev, logging
