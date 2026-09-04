package main

import (
	"log/slog"

	"github.com/biplob-codes/capto/internal/db"
)

type application struct {
	logger *slog.Logger
	db     *db.Queries
}
