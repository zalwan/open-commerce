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
PID="$(echo "$PRODUCTS" | python3 -c 'import json,sys; print(json.load(sys.stdin)["items"][0]["id"])')"

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

echo "== catalog: categories + pagination =="
curl -fsS "$BASE/api/v1/categories" | grep -q Fashion || fail "categories"
curl -fsS "$BASE/api/v1/products?category=Fashion&sort=price_desc&per_page=1&page=1" | grep -q '"total":' || fail "pagination"

echo "== cart set-qty =="
curl -fsS -X POST "$BASE/api/v1/cart/items" -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"$SESS-q\",\"product_id\":\"$PID\",\"qty\":1}" -o /dev/null || fail "cart add for setqty"
curl -fsS -X PUT "$BASE/api/v1/cart/items" -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"$SESS-q\",\"product_id\":\"$PID\",\"qty\":2}" | grep -q '"qty":2' || fail "set-qty"

echo "== order tracking =="
curl -fsS "$BASE/api/v1/orders?email=smoke@test.local" | grep -q "$OID" || fail "tracking by email"

echo "== admin stats =="
curl -fsS "$BASE/api/v1/admin/stats" -H "Authorization: Bearer $TOKEN" | grep -q '"products":' || fail "stats"

echo "== image upload =="
python3 - "$BASE" "$TOKEN" <<'EOF' || fail "upload"
import json, sys, urllib.request
base, token = sys.argv[1], sys.argv[2]
png = bytes.fromhex('89504e470d0a1a0a0000000d4948445200000001000000010802000000907753de0000000c4944415408d763f8ffff3f0005fe02fedccc59e70000000049454e44ae426082')
def api(method, path, body=None, ctype='application/json', auth=True):
    h = {}
    if auth: h['Authorization'] = f'Bearer {token}'
    data = json.dumps(body).encode() if isinstance(body, dict) else body
    if ctype: h['Content-Type'] = ctype
    r = urllib.request.Request(base + path, method=method, data=data, headers=h)
    return urllib.request.urlopen(r)
api('POST', '/api/v1/admin/products', {'id': 'p-smoke-img', 'name': 'SmokeImg', 'priceMinor': 1000, 'stock': 2})
b = ('--B\r\nContent-Disposition: form-data; name="image"; filename="t.png"\r\nContent-Type: image/png\r\n\r\n').encode() + png + b'\r\n--B--\r\n'
p = json.load(api('POST', '/api/v1/admin/products/p-smoke-img/image', b, 'multipart/form-data; boundary=B'))
assert p['imageUrl'] == '/static/p-smoke-img.png', p
img = api('GET', p['imageUrl'], ctype=None, auth=False)
assert img.status == 200, img.status
api('DELETE', '/api/v1/admin/products/p-smoke-img')
print('upload ok')
EOF

echo "SMOKE OK"
