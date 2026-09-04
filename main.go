package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found, falling back the system variable")
	}
	ctx := context.Background()
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		logger.Error("database url not provided")
		os.Exit(1)
	}
	logger.Info("connecting to database")
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		logger.Error("failed to create pool", "error", err)
		os.Exit(1)
	}
	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to connect to db", "error", err)
		os.Exit(1)

	}
	defer pool.Close()
	dbPool := db.New(pool)
	validate := validator.New()
	app := application{logger: logger, db: dbPool, validate: validate}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      app.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("server starting", "port", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
	logger.Warn("server stopped")
}
