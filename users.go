package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sudhanshu069/chirpy/internal/auth"
	"github.com/Sudhanshu069/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := payload{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedIn, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedIn,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	req := payload{}

	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Email and password are required",
		)
		return
	}

	if req.ExpiresInSeconds <= 0 || req.ExpiresInSeconds > 3600 {
		req.ExpiresInSeconds = 3600
	}

	user, err := cfg.db.LoginUser(r.Context(), req.Email)
	if err != nil {
		log.Print(err)
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
		)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)

	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if match {
		token, err := auth.MakeJWT(user.ID, cfg.jwtsecret, time.Second*time.Duration(req.ExpiresInSeconds))
		if err != nil {
			log.Print(err)
			respondWithError(
				w,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		})
	} else {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
		)
		return
	}

}
