package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gcancel/chirpy/internal/auth"
	"github.com/gcancel/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsCreate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type:", "application/json")

	type parameters struct {
		Body string `json:"body"`
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

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error retrieving bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid JWT token", err)
		return
	}
	fmt.Printf("Token: %v\n", userID)

	cleanedBody := cleanChirp(response.Body)

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error saving chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type:", "application/json")

	authorParam := req.URL.Query().Get("author_id")
	sortParam := req.URL.Query().Get("sort")
	if authorParam != "" {
		authorID, err := uuid.Parse(authorParam)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "error parsing author ID", err)
			return
		}

		chirps, err := cfg.dbQueries.GetChirpByUser(req.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error retrieving chirps from database", err)
			return
		}
		chirpResults := createChirpSlice(chirps, sortParam)
		respondWithJSON(w, http.StatusOK, chirpResults)
		return
	} else {
		chirps, err := cfg.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error retrieving chirps", err)
			return
		}
		chirpResults := createChirpSlice(chirps, sortParam)
		respondWithJSON(w, http.StatusOK, chirpResults)
		return
	}

}

func (cfg *apiConfig) handleGetUserChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type:", "application/json")

	chirpParam := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing uuid parameter", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error retrieving chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func createChirpSlice(chirps []database.Chirp, order string) []Chirp {
	chirpSlice := make([]Chirp, 0)
	for _, chirp := range chirps {
		chirpSlice = append(chirpSlice, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	if order == "desc" {
		sort.Slice(chirpSlice, func(i, j int) bool { return chirpSlice[i].CreatedAt.After(chirpSlice[j].CreatedAt) })
	}
	return chirpSlice
}
