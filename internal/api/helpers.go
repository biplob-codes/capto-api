package api

import (
	"encoding/json"
	"net/http"
)

func (a *Application) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		a.Logger.Error("failed to encode response", "error", err)
	}
}

func (a *Application) errorResponse(w http.ResponseWriter, status int, message string) {
	a.writeJSON(w, status, map[string]string{"error": message})
}

func (a *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.Logger.Warn("bad request", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusBadRequest, "invalid request")
}

func (a *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.Logger.Warn("validation failed", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
}

func (a *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.Logger.Error("internal server error", "path", r.URL.Path, "error", err)
	a.errorResponse(w, http.StatusInternalServerError, "something went wrong")
}

func (a *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	a.errorResponse(w, http.StatusNotFound, "resource not found")
}

func (a *Application) noContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
