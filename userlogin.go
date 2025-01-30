package main

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/NickLiu-0717/Chirpy/internal/auth"
	"github.com/NickLiu-0717/Chirpy/internal/database"
)

func (cfg *apiConfig) userlogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var usr userInfo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		log.Printf("Error decoding json: %s", err)
		w.WriteHeader(500)
		return
	}

	dbUser, err := cfg.db.GetUser(r.Context(), usr.Email)
	if err != nil {
		if err = respondWithError(w, 401, "Incorrect email or password"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	err = auth.CheckPasswordHash(usr.Password, dbUser.HashedPassword)
	if err != nil {
		if err = respondWithError(w, 401, "Incorrect email or password"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	token, err := auth.MakeJWT(dbUser.ID, cfg.secretKey)
	if err != nil {
		log.Printf("Error making token: %s", err)
		w.WriteHeader(500)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error making fresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	dbRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		log.Printf("Error adding fresh token to table: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: dbRefreshToken.Token,
	}
	if err = respondWithJSON(w, 200, user); err != nil {
		log.Printf("Error responding with json: %s", err)
		w.WriteHeader(500)
		return
	}
}
