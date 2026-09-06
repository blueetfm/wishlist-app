package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"

	"github.com/blueetfm/wishlist-app/backend/internal/api/routes"
	"github.com/blueetfm/wishlist-app/backend/internal/config"
	"github.com/blueetfm/wishlist-app/backend/internal/db"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{cfg.SupabaseJWKSURL})
	if err != nil {
		log.Fatalf("failed to load JWKS from %s: %v", cfg.SupabaseJWKSURL, err)
	}

	deps := routes.Dependencies{
		WishlistService: services.NewWishlistService(pool),
		JWKS:            jwks,
	}

	mux := routes.NewRouter(deps)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("wishlist-api listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
