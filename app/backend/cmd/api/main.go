// Command api runs the Open Commerce backend.
package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/open-commerce/backend/internal/auth"
	"github.com/open-commerce/backend/internal/cart"
	"github.com/open-commerce/backend/internal/catalog"
	httpapi "github.com/open-commerce/backend/internal/http"
	"github.com/open-commerce/backend/internal/order"
	"github.com/open-commerce/backend/internal/payment"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	catalogSvc := catalog.NewService(catalog.NewMemoryStore())
	cartSvc := cart.NewService(catalogSvc)
	pay := payment.NewStub()
	orderSvc := order.NewService(cartSvc, pay)
	authSvc := auth.NewService()

	router := httpapi.NewRouter(httpapi.Deps{
		Catalog: catalogSvc,
		Cart:    cartSvc,
		Order:   orderSvc,
		Auth:    authSvc,
		Pay:     pay,
	}, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	logger.Info("open-commerce api listening", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
