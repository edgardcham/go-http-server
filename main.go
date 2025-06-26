package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	// This defines a multiplexer/router
	mux := http.NewServeMux()
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// this instantiates the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// we wrap the handler with the middlewareMetricsInc function
	mux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerResetMetrics)
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server.ListenAndServe()
}
