package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gcancel/chirpy/internal/auth"
	"github.com/gcancel/chirpy/internal/database"
)

func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string
		Password string
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var response parameters
	err := decoder.Decode(&response)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding JSON", err)
		return
	}

	hashedPassword, err := auth.HashPassword(response.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error storing password", err)
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          response.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new user.", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
	fmt.Printf("User added to database: %v", user.Email)

}
