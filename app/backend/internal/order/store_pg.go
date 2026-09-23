package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/open-commerce/backend/internal/cart"
)

// PGStore persists orders in Postgres. Items are JSONB snapshots (ADR-0003).
type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore { return &PGStore{pool: pool} }

func (s *PGStore) NextID(ctx context.Context) (string, error) {
	var n int64
	if err := s.pool.QueryRow(ctx, `SELECT nextval('order_seq')`).Scan(&n); err != nil {
		return "", fmt.Errorf("next order id: %w", err)
	}
	return fmt.Sprintf("o-%d", n), nil
}

func (s *PGStore) Save(ctx context.Context, o Order) error {
	raw, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("encode order items: %w", err)
	}
	if o.Items == nil {
		raw = []byte("[]")
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO orders (id, session_id, email, items, total_minor, currency, status, payment_tx, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status`,
		o.ID, o.SessionID, o.Email, raw, o.TotalMinor, o.Currency, string(o.Status), o.PaymentTx, o.CreatedAt)
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}
	return nil
}

func (s *PGStore) Get(ctx context.Context, id string) (Order, error) {
	var o Order
	var raw []byte
	var status string
	err := s.pool.QueryRow(ctx, `SELECT id, session_id, email, items, total_minor, currency, status, payment_tx, created_at FROM orders WHERE id = $1`, id).
		Scan(&o.ID, &o.SessionID, &o.Email, &raw, &o.TotalMinor, &o.Currency, &status, &o.PaymentTx, &o.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, ErrNotFound
	}
	if err != nil {
		return Order{}, fmt.Errorf("get order: %w", err)
	}
	o.Status = Status(status)
	if len(raw) > 0 {
		o.Items = []cart.Item{}
		if err := json.Unmarshal(raw, &o.Items); err != nil {
			return Order{}, fmt.Errorf("decode order items: %w", err)
		}
	} else {
		o.Items = []cart.Item{}
	}
	return o, nil
}

func (s *PGStore) List(ctx context.Context) ([]Order, error) {
	return s.listWhere(ctx, "")
}

func (s *PGStore) ListByEmail(ctx context.Context, email string) ([]Order, error) {
	return s.listWhere(ctx, email)
}

func (s *PGStore) listWhere(ctx context.Context, email string) ([]Order, error) {
	query := `SELECT id, session_id, email, items, total_minor, currency, status, payment_tx, created_at FROM orders ORDER BY created_at DESC`
	args := []any{}
	if email != "" {
		query = `SELECT id, session_id, email, items, total_minor, currency, status, payment_tx, created_at FROM orders WHERE email = $1 ORDER BY created_at DESC`
		args = append(args, email)
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()
	out := []Order{}
	for rows.Next() {
		var o Order
		var raw []byte
		var status string
		if err := rows.Scan(&o.ID, &o.SessionID, &o.Email, &raw, &o.TotalMinor, &o.Currency, &status, &o.PaymentTx, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		o.Status = Status(status)
		o.Items = []cart.Item{}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &o.Items); err != nil {
				return nil, fmt.Errorf("decode order items: %w", err)
			}
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
