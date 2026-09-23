#!/bin/sh
# Smoke test E2E minimal. Requires API running at $BASE (default localhost:8080).
set -eu
BASE="${BASE:-http://localhost:8080}"

fail() { echo "SMOKE FAIL: $1" >&2; exit 1; }

echo "== healthz =="
curl -fsS "$BASE/healthz" | grep -q ok || fail "healthz"

echo "== products =="
PRODUCTS="$(curl -fsS "$BASE/api/v1/products")"
echo "$PRODUCTS" | grep -q "p-kaos-hitam" || fail "seed products missing"
PID="$(echo "$PRODUCTS" | python3 -c 'import json,sys; print(json.load(sys.stdin)[0]["id"])')"

echo "== cart + checkout =="
SESS="smoke-$(date +%s)"
curl -fsS -X POST "$BASE/api/v1/cart/items" -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"$SESS\",\"product_id\":\"$PID\",\"qty\":1}" | grep -q subtotalMinor || fail "add cart"
ORDER="$(curl -fsS -X POST "$BASE/api/v1/orders/checkout" -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"$SESS\",\"email\":\"smoke@test.local\"}")"
echo "$ORDER" | grep -q '"status":"paid"' || fail "checkout not paid: $ORDER"

echo "SMOKE OK"
