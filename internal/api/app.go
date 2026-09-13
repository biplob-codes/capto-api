package api

import (
	"log/slog"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/go-playground/validator/v10"
)

type Application struct {
	Logger   *slog.Logger
	Db       *db.Queries
	Validate *validator.Validate
}
