package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)
func healthHandler(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Server healthy!"))
}
func main(){
    logger:=slog.New(slog.NewJSONHandler(os.Stdout,nil))

	mux:=http.NewServeMux()

	mux.HandleFunc("GET /health" ,healthHandler)

    server:=&http.Server{
		Addr:":8080",
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
		IdleTimeout: 60*time.Second,
		Handler: mux,
	}

	logger.Info("server starting","port",server.Addr)

	if err:=server.ListenAndServe();err!=nil && !errors.Is(err,http.ErrServerClosed){
		logger.Error("server failed","error",err)
		os.Exit(1)
	}
    logger.Warn("server stopped")
 }