package main

import "net/http"

func (a *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server healthy!"))
}
