package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/biplob-codes/capto/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RequestNote struct {
	Note string `json:"note" validate:"required,max=256"`
}

func (app *application) createRequest(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	start := time.Now()
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	reqUrl := scheme + "://" + r.Host + r.RequestURI
	remoteAddr := r.RemoteAddr
	method := r.Method
	headersByte, err := json.Marshal(r.Header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	queryParamsByte, err := json.Marshal(r.URL.Query())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	endpoint, err := app.db.GetEndpointByToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	if endpoint.Status.EndpointStatus != "ACTIVE" {
		app.errorResponse(w, http.StatusGone, "This webhook endpoint is inactive")
		return
	}
	queryParams := string(queryParamsByte)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		var errMaxRead *http.MaxBytesError
		if errors.As(err, &errMaxRead) {
			app.errorResponse(w, http.StatusRequestEntityTooLarge, "request body too large")
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	body := string(bodyByte)
	size := int32(len(bodyByte))
	headers := string(headersByte)
	pgHeaders := pgtype.Text{String: headers, Valid: len(headers) != 0}
	pgQueryParams := pgtype.Text{String: queryParams, Valid: len(queryParams) != 0}
	pgBody := pgtype.Text{String: body, Valid: len(body) != 0}

	duration := time.Since(start).Microseconds()

	req, err := app.db.CreateRequest(r.Context(), db.CreateRequestParams{Url: reqUrl, RemoteAddr: remoteAddr, BodySize: size, Method: db.ReqMethod(method), Headers: pgHeaders, EndpointID: endpoint.ID, QueryParams: pgQueryParams, Body: pgBody, Duration: int32(duration)})

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, req)
}

func (app *application) getEndpointRequests(w http.ResponseWriter, r *http.Request) {
	endpointIdParam := r.PathValue("endpointId")
	var endpointId pgtype.UUID
	if err := endpointId.Scan(endpointIdParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	req, err := app.db.GetRequestsByEndpointId(r.Context(), endpointId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, req)
}

func (app *application) getRequest(w http.ResponseWriter, r *http.Request) {
	idParams := r.PathValue("id")
	var reqId pgtype.UUID
	if err := reqId.Scan(idParams); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	req, err := app.db.GetRequestById(r.Context(), reqId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, req)
}

func (app *application) addNoteToRequests(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	var reqId pgtype.UUID
	if err := reqId.Scan(idParam); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var reqNote RequestNote
	if err := json.NewDecoder(r.Body).Decode(&reqNote); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validate.Struct(reqNote); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}
	pgNote := pgtype.Text{String: reqNote.Note, Valid: true}
	req, err := app.db.AddRequestNote(r.Context(), db.AddRequestNoteParams{Note: pgNote, ID: reqId})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, req)

}
