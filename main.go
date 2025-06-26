package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync/atomic"

	"os"

	"github.com/edgardcham/go-http-server/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	// get db url
	dbURL := os.Getenv("DB_URL")
	// open connection
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	if err != nil {
		fmt.Println("Couldn't establish connection to DB")
		os.Exit(1)
	}
	// This defines a multiplexer/router
	mux := http.NewServeMux()
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
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
