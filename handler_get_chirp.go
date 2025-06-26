package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Couldn't parse UUID")
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Couldn't find chirp")
		return
	}

	chirpResponse := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, 200, chirpResponse)
}
