package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/biplob-codes/capto/internal/utils"
	"github.com/jackc/pgx/v5"
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

func getEndpoint(e db.Endpoint) EndpointResponse {
	return EndpointResponse{
		ID:        e.ID,
		Label:     e.Label,
		Token:     e.Token,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Status:    e.Status.EndpointStatus,
		ReqCount:  e.ReqCount,
	}
}
func (app *Application) createEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpointReq CreateEndpointReq
	if err := json.NewDecoder(r.Body).Decode(&endpointReq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.Validate.Struct(endpointReq); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}
	token, err := utils.GenerateToken(21)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	endpoint, err := app.Db.CreateEndpoint(r.Context(), db.CreateEndpointParams{Label: endpointReq.Label, Token: token})
	if err != nil {
		if db.IsUniqueViolation(err) {
			app.errorResponse(w, http.StatusConflict, "token already exists")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, getEndpoint(db.Endpoint(endpoint)))
}

func (app *Application) listEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := app.Db.ListEndpoints(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var result []EndpointResponse
	for _, e := range endpoints {
		result = append(result, getEndpoint(e))
	}
	app.writeJSON(w, http.StatusOK, result)

}

func (app *Application) getEndpoint(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var id pgtype.UUID
	if err := id.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	endpoint, err := app.Db.GetEndpoint(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, getEndpoint(endpoint))
}

func (app *Application) updateEndpoint(w http.ResponseWriter, r *http.Request) {
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
	if err := app.Validate.Struct(upEnReq); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	endpoint, err := app.Db.UpdateEndpoint(r.Context(), db.UpdateEndpointParams{Status: db.NullEndpointStatus{EndpointStatus: upEnReq.Status, Valid: true}, Label: upEnReq.Label, ID: endpointId})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, getEndpoint(endpoint))

}

func (app *Application) deleteEndpoint(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var eId pgtype.UUID
	if err := eId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.Db.DeleteEndpoint(r.Context(), eId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.noContentResponse(w)
}
