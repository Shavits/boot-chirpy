package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleaned_words := []string{}
	words := strings.Split(params.Body, " ")
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range(words){
		if slices.Contains(forbiddenWords, strings.ToLower(word)){
			cleaned_words = append(cleaned_words, "****")
			continue
		}
		cleaned_words = append(cleaned_words, word)
	} 

	cleanedBody := strings.Join(cleaned_words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
}
