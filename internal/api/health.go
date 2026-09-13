package api

import "net/http"

func (a *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server healthy!"))
}
