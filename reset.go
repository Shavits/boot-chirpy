package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev"{
		respondWithError(w,403, "action not allowed", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil{
		respondWithError(w, 500, "unable to reset users", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and users reset"))

}
