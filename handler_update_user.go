package main

import (
	"encoding/json"
	"net/http"

	"github.com/edgardcham/go-http-server/internal/auth"
	"github.com/edgardcham/go-http-server/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. get jwt from header
	// 2. unmarshal the request
	// 3. Get the user corresponding to the JWT
	// 4. Hash the password and update the user's email and password
	// 5. return the User (without hashedPassword)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing JWT in Header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, 401, "Invalid JWT")
		return
	}

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	req := request{}

	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Error decoding request body")
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, "Could not find user.")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 400, "Could not hash password")
		return
	}

	userDBParams := database.UpdateUserEmailAndPassParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	newUser, err := cfg.db.UpdateUserEmailAndPass(r.Context(), userDBParams)
	if err != nil {
		respondWithError(w, 400, "Could not update user")
		return
	}

	userResponse := User{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	}
	respondWithJSON(w, 200, userResponse)

}
