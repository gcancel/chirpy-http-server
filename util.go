package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling data... %v\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		fmt.Printf("Responding with 5XX error... %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

func cleanChirp(msg string) string {

	notAllowed := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(msg, " ")
	for idx, word := range words {
		for _, naughtyWord := range notAllowed {
			if strings.ToLower(word) == naughtyWord {
				words[idx] = "****"
			}
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}
