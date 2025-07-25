package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gcancel/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaHook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	type response struct {
		// nothing
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error retrieving apikey", err)
		return
	}
	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "invalid apikey", err)
		return
	}

	var request parameters
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding json body", err)
		return
	}
	fmt.Println(request.Event)
	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing user ID", err)
		return
	}
	if request.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, response{})
	} else {
		err := cfg.dbQueries.UpdateChirpyRedStatus(req.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "error updating chirpy red status", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, response{})
	}

}
