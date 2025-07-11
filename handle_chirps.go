package main

import (
	"encoding/json"
	"net/http"
)

func handleChirp(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type:", "application/json")

	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var response parameters
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json.", err)
		return
	}

	if len(response.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long.", err)
		return
	}

	cleanedBody := cleanChirp(response.Body)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})

}
