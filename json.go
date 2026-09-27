package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorz struct {
		Error string `json:"error"`
	}

	resBody := errorz{
		Error: msg,
	}

	dat, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}

func cleanProfanity(body string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(body, " ")

	for i, word := range words {
		if profane[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func handleChirp(w http.ResponseWriter, r *http.Request) {
	type rezponse struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	res := rezponse{}
	err := decoder.Decode(&res)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	type successResponse struct {
		Message string `json:"cleaned_body"`
	}

	res.Body = cleanProfanity(res.Body)

	respondWithJSON(w, http.StatusOK, successResponse{
		Message: res.Body,
	})
}
