package main

import (
	"encoding/json"
	"fmt"
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorVals struct {
		Error string `json:"error"`
	}

	respBody := errorVals{
		Error: msg,
	}

	// marshal it to JSON
	data, err := json.Marshal(&respBody)
	if err != nil {
		// if error, log and return a 500
		fmt.Println("Was unable to send back error response")
		w.WriteHeader(500)
	}
	// otherwise return a 400 for the original error with the data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(&payload)
	if err != nil {
		fmt.Println("Was unable to send back error response")
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
