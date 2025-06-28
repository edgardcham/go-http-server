package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync/atomic"

	"os"

	"github.com/edgardcham/go-http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found")
		os.Exit(1)
	}
	// get db url
	dbURL := os.Getenv("DB_URL")
	// open connection
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	if err != nil {
		fmt.Println("Couldn't establish connection to DB")
		os.Exit(1)
	}

	defer db.Close()
	// This defines a multiplexer/router
	mux := http.NewServeMux()
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("JWT_SECRET_KEY"),
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
	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirpByID)
	mux.HandleFunc("POST /api/login", apiConfig.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiConfig.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiConfig.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handlerDeleteChirp)

	server.ListenAndServe()
}
