package main

import "net/http"

func (a *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.healthHandler)
	return mux
}
