package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-chirpy/internal/auth"
	"github.com/shavits/boot-chirpy/internal/database"
)



type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}



func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}


	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Bearer", err)
		return
	}

	userMatch, err := auth.ValidateJWT(token, cfg.secret_key)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Invalid Token", err)
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


	chirpParams := database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body: cleanedBody,
		UserID: userMatch,


	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	})
}


func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	chirps := make([]Chirp, 0, len(dbChirps))
    for _, c := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        c.ID,
            CreatedAt: c.CreatedAt,
            UpdatedAt: c.UpdatedAt,
            Body:      c.Body,
            UserId:    c.UserID,
        })
    }

    respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpIdStr)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse ID", err)
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), id)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Couldn't find Chirp", err)
		return
	}

	chirp := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}
	
	respondWithJSON(w, http.StatusOK, chirp)

}


func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpIdStr)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Bearer", err)
		return
	}

	userMatch, err := auth.ValidateJWT(token, cfg.secret_key)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Invalid Token", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), id)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userMatch{
		respondWithError(w, http.StatusForbidden, "Chirp belongs to different user", err)
		return
	}


	err = cfg.dbQueries.DeleteChirpById(r.Context(), id)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
	
}
