package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edgardcham/go-http-server/internal/auth"
	"github.com/edgardcham/go-http-server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "User not found")
		return
	}

	// check password
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// create token

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, 400, "Error creating JWT")
		return
	}

	// refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 400, "Error creating refresh token")
		return
	}

	refreshTokenDBParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshTokenDBParams)
	if err != nil {
		respondWithError(w, 400, "Couldn't store refresh token")
		return
	}

	userResponse := User{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        jwtToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, 200, userResponse)
}
