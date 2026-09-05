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
