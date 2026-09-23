package catalog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGStore persists products in Postgres. Used when DATABASE_URL is set.
type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore { return &PGStore{pool: pool} }

func (s *PGStore) List(ctx context.Context) ([]Product, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, description, category, price_minor, currency, stock, image_url, created_at FROM products ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.PriceMinor, &p.Currency, &p.Stock, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *PGStore) Get(ctx context.Context, id string) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx, `SELECT id, name, description, category, price_minor, currency, stock, image_url, created_at FROM products WHERE id = $1`, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.PriceMinor, &p.Currency, &p.Stock, &p.ImageURL, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (s *PGStore) Save(ctx context.Context, p Product) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO products (id, name, description, category, price_minor, currency, stock, image_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, description = EXCLUDED.description, category = EXCLUDED.category,
			price_minor = EXCLUDED.price_minor, currency = EXCLUDED.currency,
			stock = EXCLUDED.stock, image_url = EXCLUDED.image_url`,
		p.ID, p.Name, p.Description, p.Category, p.PriceMinor, p.Currency, p.Stock, p.ImageURL, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("save product: %w", err)
	}
	return nil
}

func (s *PGStore) Delete(ctx context.Context, id string) error {
	res, err := s.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DecrementStock is a single atomic UPDATE: concurrent checkouts cannot oversell.
func (s *PGStore) DecrementStock(ctx context.Context, id string, qty int) error {
	res, err := s.pool.Exec(ctx, `UPDATE products SET stock = stock - $2 WHERE id = $1 AND stock >= $2`, id, qty)
	if err != nil {
		return fmt.Errorf("decrement stock: %w", err)
	}
	if res.RowsAffected() == 1 {
		return nil
	}
	if _, gerr := s.Get(ctx, id); gerr != nil {
		return gerr // ErrNotFound
	}
	return ErrInsufficientStock
}

func (s *PGStore) IncrementStock(ctx context.Context, id string, qty int) error {
	res, err := s.pool.Exec(ctx, `UPDATE products SET stock = stock + $2 WHERE id = $1`, id, qty)
	if err != nil {
		return fmt.Errorf("increment stock: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SeedIfEmpty inserts the v0.1 seed catalog once, so a fresh DB matches memory dev.
func (s *PGStore) SeedIfEmpty(ctx context.Context) error {
	var count int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&count); err != nil {
		return fmt.Errorf("count products: %w", err)
	}
	if count > 0 {
		return nil
	}
	now := time.Now().UTC()
	for _, p := range SeedProducts(now) {
		p.CreatedAt = now
		if err := s.Save(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
