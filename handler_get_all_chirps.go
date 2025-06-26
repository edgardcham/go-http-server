package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Couldn't fetch chirps")
		return
	}

	chirpResponses := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpResponses[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, 200, chirpResponses)
}
