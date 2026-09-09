package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type ResponseConfig struct {
	Method  db.ResMethod `json:"method" validate:"oneof=ALL GET POST PUT PATCH DELETE OPTIONS"`
	Status  int          `json:"statusCode" validate:"omitempty,min=100,max=599"`
	Headers string       `json:"headers"`
	Body    string       `json:"body"`
	Delay   int          `json:"delay" validate:"min=0,max=30"`
}

func (app *application) createResponseConfig(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var endpointId pgtype.UUID
	if err := endpointId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var body ResponseConfig
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validate.Struct(body); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}
	pgStatus := pgtype.Int4{Int32: int32(body.Status)}
	pgDelay := pgtype.Int4{Int32: int32(body.Delay)}
	pgHeaders := pgtype.Text{String: body.Headers, Valid: len(body.Headers) > 0}
	pgBody := pgtype.Text{String: body.Body, Valid: len(body.Body) > 0}
	conf, err := app.db.CreateResponseConfig(r.Context(), db.CreateResponseConfigParams{Method: body.Method, EndpointID: endpointId, StatusCode: pgStatus, Delay: pgDelay, Headers: pgHeaders, Body: pgBody})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			app.errorResponse(w, http.StatusConflict, "config for this method already exists")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, conf)

}

func (app *application) listResponseConfigsByEndpointId(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var endpointId pgtype.UUID
	if err := endpointId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	confs, err := app.db.GetResponseConfigsByEndpointId(r.Context(), endpointId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, confs)
}

func (app *application) updateResConfig(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var resId pgtype.UUID
	if err := resId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var body ResponseConfig
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validate.Struct(body); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}
	pgStatus := pgtype.Int4{Int32: int32(body.Status), Valid: true}
	pgDelay := pgtype.Int4{Int32: int32(body.Delay), Valid: true}
	pgHeaders := pgtype.Text{String: body.Headers, Valid: len(body.Headers) > 0}
	pgBody := pgtype.Text{String: body.Body, Valid: len(body.Body) > 0}
	conf, err := app.db.UpdateResponseConfig(r.Context(), db.UpdateResponseConfigParams{Method: body.Method, StatusCode: pgStatus, Headers: pgHeaders, Body: pgBody, Delay: pgDelay, ID: resId})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			app.errorResponse(w, http.StatusConflict, "config for this method already exists")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, conf)
}

func (app *application) deleteResConfig(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var resId pgtype.UUID
	if err := resId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.db.DeleteResponseConfig(r.Context(), resId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusNoContent, "")
}
