package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/biplob-codes/capto/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEndpointReq struct {
	Label string `json:"label" validate:"required,min=3"`
}
type UpdateEndpointReq struct {
	Label  string            `json:"label" validate:"required,min=3"`
	Status db.EndpointStatus `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}
type EndpointResponse struct {
	ID        pgtype.UUID        `json:"id"`
	Label     string             `json:"label"`
	Token     string             `json:"token"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	Status    db.EndpointStatus  `json:"status"`
	ReqCount  pgtype.Int4        `json:"reqCount"`
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
	result := EndpointResponse{
		ID:        endpoint.ID,
		Label:     endpoint.Label,
		Token:     endpoint.Token,
		CreatedAt: endpoint.CreatedAt,
		UpdatedAt: endpoint.UpdatedAt,
		Status:    endpoint.Status.EndpointStatus,
		ReqCount:  endpoint.ReqCount,
	}
	app.writeJSON(w, http.StatusCreated, result)
}

func (app *application) listEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := app.db.ListEndpoints(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var result []EndpointResponse
	for _, e := range endpoints {
		r := EndpointResponse{
			ID:        e.ID,
			Label:     e.Label,
			Token:     e.Token,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			Status:    e.Status.EndpointStatus,
			ReqCount:  e.ReqCount,
		}
		result = append(result, r)
	}
	app.writeJSON(w, http.StatusOK, result)

}

func (app *application) getEndpoint(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	endpoint, err := app.db.GetEndpoint(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	result := EndpointResponse{
		ID:        endpoint.ID,
		Label:     endpoint.Label,
		Token:     endpoint.Token,
		CreatedAt: endpoint.CreatedAt,
		UpdatedAt: endpoint.UpdatedAt,
		Status:    endpoint.Status.EndpointStatus,
		ReqCount:  endpoint.ReqCount,
	}
	app.writeJSON(w, http.StatusOK, result)
}

func (app *application) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var endpointId pgtype.UUID
	if err := endpointId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var upEnReq UpdateEndpointReq
	if err := json.NewDecoder(r.Body).Decode(&upEnReq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validate.Struct(upEnReq); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	endpoint, err := app.db.UpdateEndpoint(r.Context(), db.UpdateEndpointParams{Status: db.NullEndpointStatus{EndpointStatus: upEnReq.Status, Valid: true}, Label: upEnReq.Label, ID: endpointId})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	result := EndpointResponse{
		ID:        endpoint.ID,
		Label:     endpoint.Label,
		Token:     endpoint.Token,
		CreatedAt: endpoint.CreatedAt,
		UpdatedAt: endpoint.UpdatedAt,
		Status:    endpoint.Status.EndpointStatus,
		ReqCount:  endpoint.ReqCount,
	}
	app.writeJSON(w, http.StatusOK, result)

}

func (app *application) deleteEndpoint(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var eId pgtype.UUID
	if err := eId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.db.DeleteEndpoint(r.Context(), eId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.noContentResponse(w)
}
