package main

import (
	"net/http"
	"sort"

	"github.com/edgardcham/go-http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// optional query param
	authorId := r.URL.Query().Get("author_id")
	s := r.URL.Query().Get("sort")

	// define here, because when using := it creates them within the scope.
	var chirps []database.Chirp
	var err error

	if authorId != "" {
		userId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, 400, "Couldn't parse author ID")
			return
		}

		chirps, err = cfg.db.GetAllChirpsForUser(r.Context(), userId) // if chirps, err := then the inner chirp will shadow the outer one, hence it's undesired
		if err != nil {
			respondWithError(w, 400, "Couldn't fetch chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 400, "Couldn't fetch chirps")
			return
		}
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

	if s == "" || s == "asc" {
		respondWithJSON(w, 200, chirpResponses)
		return
	} else {
		sort.Slice(chirpResponses, func(i, j int) bool {
			return chirpResponses[i].CreatedAt.After(chirpResponses[j].CreatedAt)
		})
	}

	respondWithJSON(w, 200, chirpResponses)
}
