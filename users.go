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
		Password string `json:"password"`
		Email    string `json:"email"`
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
		token, err := auth.MakeJWT(user.ID, cfg.jwtsecret, time.Hour)
		if err != nil {
			log.Print(err)
			respondWithError(
				w,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
			return
		}

		refreshToken := auth.MakeRefreshToken()

		rfToken, err := cfg.db.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
		})

		respondWithJSON(w, http.StatusOK, User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: rfToken.Token,
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

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Res struct {
		Body string `json:"token"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	val, err := cfg.db.SearchRefreshToken(r.Context(), tokenString)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not matching")
		return
	}

	if val.ExpiresAt.Before(time.Now().UTC()) || val.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "not matching")
		return
	}

	newToken, err := auth.MakeJWT(val.UserID, cfg.jwtsecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating new token")
	}

	respondWithJSON(w, http.StatusOK, Res{
		Body: newToken,
	})

}
