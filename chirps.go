package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
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

	// _, err := cfg.dbQueries.GetUser(context.Background(), userName)
	// if err == sql.ErrNoRows{
	// 	return fmt.Errorf("user %s does not exist", userName)
	// }

	
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
		UserID: params.UserId,


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
