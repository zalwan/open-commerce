// Package db opens the Postgres pool and applies embedded migrations.
// Used only when DATABASE_URL is set; otherwise the API runs on memory stores.
package db

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// migrations holds DDL applied in lexical order at startup.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Open creates a bounded pool. Callers Ping to verify before serving.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}

// Migrate executes each embedded *.sql file statement-by-statement.
// Files contain plain DDL without semicolons inside literals, so a
// semicolon split is sufficient (no external migrate tool in v0.2).
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		raw, err := migrations.ReadFile("migrations/" + e.Name())
		if err != nil {
			return fmt.Errorf("read %s: %w", e.Name(), err)
		}
		// Strip full-line comments before splitting so semicolons
		// inside comments can never become statement fragments.
		lines := strings.Split(string(raw), "\n")
		nocomment := make([]string, 0, len(lines))
		for _, l := range lines {
			if t := strings.TrimSpace(l); !strings.HasPrefix(t, "--") {
				nocomment = append(nocomment, l)
			}
		}
		for _, stmt := range strings.Split(strings.Join(nocomment, "\n"), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("apply %s: %w", e.Name(), err)
			}
		}
	}
	return nil
}
