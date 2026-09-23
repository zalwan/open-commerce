-- Open Commerce v0.2 baseline schema (additive only, no destructive DDL).

CREATE TABLE IF NOT EXISTS products (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  price_minor BIGINT NOT NULL CHECK (price_minor > 0),
  currency    TEXT NOT NULL DEFAULT 'IDR',
  stock       INT NOT NULL CHECK (stock >= 0),
  image_url   TEXT NOT NULL DEFAULT '',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products (name);

CREATE TABLE IF NOT EXISTS carts (
  session_id TEXT PRIMARY KEY,
  items      JSONB NOT NULL DEFAULT '[]',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
  id          TEXT PRIMARY KEY,
  session_id  TEXT NOT NULL,
  email       TEXT NOT NULL,
  items       JSONB NOT NULL DEFAULT '[]',
  total_minor BIGINT NOT NULL CHECK (total_minor >= 0),
  currency    TEXT NOT NULL DEFAULT 'IDR',
  status      TEXT NOT NULL,
  payment_tx  TEXT NOT NULL DEFAULT '',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_created ON orders (created_at DESC);

-- Order ID sequence; app prefixes with 'o-'.
CREATE SEQUENCE IF NOT EXISTS order_seq;
