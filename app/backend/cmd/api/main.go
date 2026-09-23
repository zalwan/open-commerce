// Command api runs the Open Commerce backend.
// Memory stores by default; Postgres when DATABASE_URL is set (migrations auto-applied).
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/open-commerce/backend/internal/auth"
	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	"github.com/open-commerce/backend/internal/db"
	httpapi "github.com/open-commerce/backend/internal/http"
	"github.com/open-commerce/backend/internal/media"
	"github.com/open-commerce/backend/internal/order"
	"github.com/open-commerce/backend/internal/payment"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var (
		catalogStore catalog.Store = catalog.NewMemoryStore()
		cartStore    cart.Store    = cart.NewMemoryStore()
		orderStore   order.Store   = order.NewMemoryStore()
		storeName                  = "memory"
		ready        func(context.Context) error
	)

	if url := os.Getenv("DATABASE_URL"); url != "" {
		pool, err := db.Open(ctx, url)
		if err != nil {
			logger.Error("db open failed", "err", err)
			os.Exit(1)
		}
		defer pool.Close()
		if err := db.Migrate(ctx, pool); err != nil {
			logger.Error("migrate failed", "err", err)
			os.Exit(1)
		}
		pgCatalog := catalog.NewPGStore(pool)
		if err := pgCatalog.SeedIfEmpty(ctx); err != nil {
			logger.Error("seed failed", "err", err)
			os.Exit(1)
		}
		catalogStore = pgCatalog
		cartStore = cart.NewPGStore(pool)
		orderStore = order.NewPGStore(pool)
		storeName = "postgres"
		ready = pool.Ping
	}

	catalogSvc := catalog.NewService(catalogStore)
	cartSvc := cart.NewService(cartStore, catalogSvc)
	pay := payment.NewStub()
	orderSvc := order.NewService(orderStore, cartSvc, catalogSvc, pay)
	authSvc := auth.NewService()

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data/uploads"
	}
	mediaStore := media.NewStore(dataDir, "/static/")

	router := httpapi.NewRouter(httpapi.Deps{
		Catalog: catalogSvc,
		Cart:    cartSvc,
		Order:   orderSvc,
		Auth:    authSvc,
		Pay:     pay,
		Media:   mediaStore,
		Ready:   ready,
	}, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	logger.Info("open-commerce api listening", "addr", addr, "store", storeName)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
