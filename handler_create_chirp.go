package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/edgardcham/go-http-server/internal/auth"
	"github.com/edgardcham/go-http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// decoder on the request body
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong.")
		return
	}

	// Get JWT and validate it
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bearer token issue")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid Bearer token")
		return
	}

	chirpBody := params.Body

	if len(chirpBody) > 140 {
		respondWithError(w, 400, "Chirp is too long.")
		return
	}

	strArr := strings.Split(chirpBody, " ")
	for i, str := range strArr {
		lowercaseStr := strings.ToLower(str)
		if lowercaseStr == "kerfuffle" || lowercaseStr == "sharbert" || lowercaseStr == "fornax" {
			strArr[i] = "****"
		}
	}
	cleanedBody := strings.Join(strArr, " ")

	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 400, "Could not create chirp")
		return
	}

	// map to Chirp Struct
	chirpResp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    userID,
	}

	respondWithJSON(w, 201, chirpResp)
}
