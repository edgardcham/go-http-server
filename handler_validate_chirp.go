package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// decoder on the request body
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong.")
	}

	chirp := params.Body

	if len(chirp) > 140 {
		respondWithError(w, 400, "Chirp is too long.")
	}

	strArr := strings.Split(chirp, " ")
	for i, str := range strArr {
		lowercaseStr := strings.ToLower(str)
		if lowercaseStr == "kerfuffle" || lowercaseStr == "sharbert" || lowercaseStr == "fornax" {
			strArr[i] = "****"
		}
	}
	cleanedBody := strings.Join(strArr, " ")

	payload := returnVals{
		CleanedBody: cleanedBody,
	}

	respondWithJSON(w, 200, payload)
}
