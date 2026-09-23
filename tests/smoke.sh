#!/bin/sh
# Smoke test E2E minimal. Requires API running at $BASE (default localhost:8080).
# Covers: health, catalog, cart, checkout, admin login + product CRUD + order status.
set -eu
BASE="${BASE:-http://localhost:8080}"

fail() { echo "SMOKE FAIL: $1" >&2; exit 1; }

echo "== healthz =="
curl -fsS "$BASE/healthz" | grep -q ok || fail "healthz"
curl -fsS "$BASE/readyz" | grep -q ready || fail "readyz"

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
OID="$(echo "$ORDER" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')"

echo "== admin =="
TOKEN="$(curl -fsS -X POST "$BASE/api/v1/auth/login" -H 'Content-Type: application/json' \
  -d '{"email":"admin@shop.test","password":"admin123"}' | python3 -c 'import json,sys; print(json.load(sys.stdin)["token"])')"
[ -n "$TOKEN" ] || fail "admin login"
curl -fsS -X POST "$BASE/api/v1/admin/products" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"id":"p-smoke","name":"Smoke","priceMinor":1000,"stock":5}' | grep -q "p-smoke" || fail "admin create"
curl -fsS -X DELETE "$BASE/api/v1/admin/products/p-smoke" -H "Authorization: Bearer $TOKEN" -o /dev/null || fail "admin delete"
curl -fsS "$BASE/api/v1/admin/orders" -H "Authorization: Bearer $TOKEN" | grep -q "$OID" || fail "admin orders"
curl -fsS -X POST "$BASE/api/v1/admin/orders/$OID/status" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"status":"shipped"}' | grep -q '"status":"shipped"' || fail "order status"

echo "== admin guard =="
if curl -fsS "$BASE/api/v1/admin/orders" >/dev/null 2>&1; then fail "admin without token should 401/403"; fi

echo "== oversell guard =="
LIM="p-limited-$(date +%s)"
curl -fsS -X POST "$BASE/api/v1/admin/products" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d "{\"id\":\"$LIM\",\"name\":\"Limited\",\"priceMinor\":1000,\"stock\":1}" -o /dev/null || fail "create limited"
curl -fsS -X POST "$BASE/api/v1/cart/items" -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"smoke-a\",\"product_id\":\"$LIM\",\"qty\":1}" -o /dev/null || fail "cart a"
curl -fsS -X POST "$BASE/api/v1/cart/items" -H 'Content-Type: application/json' \
  -d '{"session_id":"smoke-b","product_id":"'"$LIM"'","qty":1}' -o /dev/null || fail "cart b"
curl -fsS -X POST "$BASE/api/v1/orders/checkout" -H 'Content-Type: application/json' \
  -d '{"session_id":"smoke-a","email":"a@test.local"}' -o /dev/null || fail "checkout a"
if curl -fsS -X POST "$BASE/api/v1/orders/checkout" -H 'Content-Type: application/json' \
  -d '{"session_id":"smoke-b","email":"b@test.local"}' >/dev/null 2>&1; then fail "oversell checkout should 409"; fi
curl -fsS -X DELETE "$BASE/api/v1/admin/products/$LIM" -H "Authorization: Bearer $TOKEN" -o /dev/null || fail "cleanup limited"

echo "SMOKE OK"
