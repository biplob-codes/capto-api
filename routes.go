package main

import "net/http"

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.healthHandler)
	mux.HandleFunc("POST /endpoints", app.createEndpoint)
	return mux
}
