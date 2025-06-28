package main

import (
	"net/http"

	"github.com/edgardcham/go-http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "JWT not set in headers")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid JWT")
		return
	}

	// get chirp by ID, if not found, return 404
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Did not specify chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	// verify it belongs to user
	if chirp.UserID != userID {
		respondWithError(w, 403, "Unauthorized")
		return
	}

	// if it does delete and return 403
	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 400, "Could not delete chirp")
		return
	}

	respondWithJSON(w, 204, nil)
}
