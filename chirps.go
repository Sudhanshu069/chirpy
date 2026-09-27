package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sudhanshu069/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateChirps(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	req := payload{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   req.Body,
		UserID: userID,
	})

	if err != nil {
		log.Printf("create chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("cannot get list: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := []Chirp{}

	for _, item := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}
