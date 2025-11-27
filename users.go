package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shavits/boot-chirpy/internal/auth"
	"github.com/shavits/boot-chirpy/internal/database"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool 	`json:"is_chirpy_red"`
}


func (cfg *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	} 

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: params.Email,
		HashedPassword: hashedPass,
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		//ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Unable decrypt password", err)
		return
	}

	if !match{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// actualDurationInSeconds := max(params.ExpiresInSeconds, 3600)
	// duraion, err := time.ParseDuration(fmt.Sprintf("%ds",actualDurationInSeconds))
	// if err != nil{
	// 	respondWithError(w, http.StatusUnauthorized, "Invalid duration for token expiry", err)
	// 	return
	// }
	token, err := auth.MakeJWT(user.ID, cfg.secret_key, time.Hour)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unable to create JWT", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to create refreshtoken", err)
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		RevokedAt: sql.NullTime{},
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to create refreshtoken", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: token,
		RefreshToken: refreshToken,
	})

}


func (cfg *apiConfig) handlerUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
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
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	} 


	updateUserParams := database.UpdateEmailAndPwdByIdParams{
		ID: userMatch,
		Email: params.Email,
		HashedPassword: hashedPass,
	}

	updatedUser, err := cfg.dbQueries.UpdateEmailAndPwdById(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	} 


	respondWithJSON(w, http.StatusOK, User{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})

	
}
