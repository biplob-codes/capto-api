package main

import "net/http"

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.healthHandler)
	mux.HandleFunc("POST /endpoints", app.createEndpoint)
	mux.HandleFunc("GET /endpoints", app.listEndpoints)
	mux.HandleFunc("GET /endpoints/{id}", app.getEndpoint)
	mux.HandleFunc("PATCH /endpoints/{id}", app.UpdateEndpoint)
	mux.HandleFunc("GET /endpoints/{endpointId}/requests", app.getEndpointRequests)
	mux.HandleFunc("/hooks/{token}", app.createRequest)
	mux.HandleFunc("GET /requests/{id}", app.getRequest)
	mux.HandleFunc("PATCH /requests/{id}/note", app.addNoteToRequests)

	return mux
}
