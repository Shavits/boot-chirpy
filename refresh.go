package main

import (
	"net/http"
	"time"

	"github.com/shavits/boot-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct{
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Bearer", err)
		return
	}
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unable to create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{Token: accessToken})


}



func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Bearer", err)
		return
	}
	err = cfg.dbQueries.RevokeToken(r.Context(), refreshToken)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke token", err)
		return
	}
	
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(http.StatusText(http.StatusNoContent)))



}
