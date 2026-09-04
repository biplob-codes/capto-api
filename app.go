package main

import (
	"log/slog"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/go-playground/validator/v10"
)

type application struct {
	logger   *slog.Logger
	db       *db.Queries
	validate *validator.Validate
}
