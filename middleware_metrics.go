package main

import (
	"fmt"
	"net/http"
)

// middlewareMetricsInc is a middleware that increments the fileserverHits counter
// it wraps the next handler and increments the counter
// Effectively we wrap a handler and return the same handler after we execute some code
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// handlerMetrics is a handler that returns the number of hits to the fileserver
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
			</html>
	`, cfg.fileserverHits.Load())))
}

// handlerResetMetrics is a handler that resets the fileserverHits counter
func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
