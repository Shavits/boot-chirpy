package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shavits/boot-chirpy/internal/auth"
)

const UPGRADED_EVENT string = "user.upgraded"

func (cfg *apiConfig) handlerUgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldn't get apiKey", err)
		return
	}

	if apiKey != cfg.polkaKey{
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event == UPGRADED_EVENT{
		_,  err = cfg.dbQueries.UpgradeToChirpyRedById(r.Context(), params.Data.UserID)
		if err != nil{
			if err == sql.ErrNoRows{
				respondWithError(w, http.StatusNotFound, "User Not Found", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Unable to upgrade user", err)
			return
		}
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))
}