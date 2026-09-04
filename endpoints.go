package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/biplob-codes/capto/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateEndpointReq struct {
	Label string `json:"label" validate:"required,min=3"`
}

func (app *application) createEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpointReq CreateEndpointReq
	if err := json.NewDecoder(r.Body).Decode(&endpointReq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validate.Struct(endpointReq); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}
	token, err := utils.GenerateToken(21)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	endpoint, err := app.db.CreateEndpoint(r.Context(), db.CreateEndpointParams{Label: endpointReq.Label, Token: token})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			app.errorResponse(w, http.StatusConflict, "token already exists")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, endpoint)
}

func (app *application) listEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := app.db.ListEndpoints(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, endpoints)

}
