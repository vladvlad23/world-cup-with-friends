package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"worldcup/backend/internal/api"
	"worldcup/backend/internal/auth"
	"worldcup/backend/internal/db"
	"worldcup/backend/internal/seed"
)

func main() {
	dbURL := env("DATABASE_URL", "postgres://worldcup:worldcup@localhost:5432/worldcup?sslmode=disable")
	jwtSecret := env("JWT_SECRET", "dev-secret")
	port := env("PORT", "8080")

	ctx := context.Background()

	pool, err := connectWithRetry(ctx, dbURL, 15)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	if err := seed.Run(ctx, pool); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Println("database ready and seeded")

	srv := api.NewServer(pool, auth.New(jwtSecret))

	httpSrv := &http.Server{
		Addr:    ":" + port,
		Handler: srv.Router(),
	}
	log.Printf("listening on :%s", port)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func connectWithRetry(ctx context.Context, url string, attempts int) (pool *pgxpool.Pool, err error) {
	for i := 0; i < attempts; i++ {
		pool, err = db.Connect(ctx, url)
		if err == nil {
			return pool, nil
		}
		log.Printf("waiting for database (%d/%d): %v", i+1, attempts, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
