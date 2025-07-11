package main

import (
	"encoding/json"
	"net/http"

	"github.com/edgardcham/go-http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "No API Key set")
		return
	}
	if key != cfg.polkaAPIKey {
		respondWithError(w, 401, "Invalid API Key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Couldn't decode payload")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "Couldn't parse user_id")
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	respondWithJSON(w, 204, nil)
}
