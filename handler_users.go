package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var response parameters
	err := decoder.Decode(&response)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding JSON", err)
		return
	}
	ctx := context.Background()
	user, err := cfg.dbQueries.CreateUser(ctx, response.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new user.", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	})
	fmt.Printf("User added to database: %v", user.Email)

}
