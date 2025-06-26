package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	respBody := errorVals{
		Error: msg,
	}

	// marshal it to JSON
	data, err := json.Marshal(&respBody)
	if err != nil {
		// if error, log and return a 500
		fmt.Println("Was unable to send back error response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}
	// otherwise return a 400 for the original error with the data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Was unable to send back error response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
