package main

import (
	"net/http"
	"time"

	"github.com/edgardcham/go-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Couldn't get refresh token from Header")
		return
	}

	refreshTokenDB, err := cfg.db.GetRefreshTokenByID(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Invalid refresh token")
		return
	}

	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	if refreshTokenDB.RevokedAt.Valid { // valid checks if not Null
		respondWithError(w, 401, "Refresh Token Revoked")
		return
	}

	jwtToken, err := auth.MakeJWT(refreshTokenDB.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 400, "Couldn't create JWT")
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}

	respondWithJSON(w, 200, response)

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Couldn't get refresh token from Header")
		return
	}

	if err := cfg.db.RevokeTokenByID(r.Context(), refreshToken); err != nil {
		respondWithError(w, 400, "Couldn't revoke the refresh token")
		return
	}

	respondWithJSON(w, 204, nil)

}
