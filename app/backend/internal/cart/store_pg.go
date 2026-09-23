package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGStore persists carts in Postgres as JSONB snapshots.
type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore { return &PGStore{pool: pool} }

func (s *PGStore) Get(ctx context.Context, sessionID string) (Cart, error) {
	var raw []byte
	err := s.pool.QueryRow(ctx, `SELECT items FROM carts WHERE session_id = $1`, sessionID).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return Cart{}, ErrCartNotFound
	}
	if err != nil {
		return Cart{}, fmt.Errorf("get cart: %w", err)
	}
	items := []Item{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &items); err != nil {
			return Cart{}, fmt.Errorf("decode cart items: %w", err)
		}
	}
	c := Cart{SessionID: sessionID, Items: items}
	c.recalc()
	return c, nil
}

func (s *PGStore) Save(ctx context.Context, c Cart) error {
	raw, err := json.Marshal(c.Items)
	if err != nil {
		return fmt.Errorf("encode cart items: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO carts (session_id, items, updated_at)
		VALUES ($1, $2, now())
		ON CONFLICT (session_id) DO UPDATE SET items = EXCLUDED.items, updated_at = now()`,
		c.SessionID, raw)
	if err != nil {
		return fmt.Errorf("save cart: %w", err)
	}
	return nil
}

func (s *PGStore) Delete(ctx context.Context, sessionID string) error {
	if _, err := s.pool.Exec(ctx, `DELETE FROM carts WHERE session_id = $1`, sessionID); err != nil {
		return fmt.Errorf("delete cart: %w", err)
	}
	return nil
}
