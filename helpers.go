package main

import (
	"encoding/json"
	"net/http"
)

func (a *application) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		a.logger.Error("failed to encode response", "error", err)
	}
}

func (a *application) errorResponse(w http.ResponseWriter, status int, message string) {
	a.writeJSON(w, status, map[string]string{"error": message})
}

func (a *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Warn("bad request", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusBadRequest, "invalid request")
}

func (a *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Warn("validation failed", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
}

func (a *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error("internal server error", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusInternalServerError, "something went wrong")
}

func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	a.errorResponse(w, http.StatusNotFound, "resource not found")
}
