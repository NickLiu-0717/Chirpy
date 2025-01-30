package main

import (
	"log"
	"net/http"

	"github.com/NickLiu-0717/Chirpy/internal/auth"
)

func (cfg *apiConfig) checkRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err = respondWithError(w, 401, "The token is not found or expired"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	accessToken, err := auth.MakeJWT(userID, cfg.secretKey)
	if err != nil {
		log.Printf("Error making jwt: %s", err)
		w.WriteHeader(500)
		return
	}
	if err = respondWithJSON(w, 200, AccessToken{Token: accessToken}); err != nil {
		log.Printf("Error responding json: %s", err)
		w.WriteHeader(500)
		return
	}
}

func (cfg *apiConfig) revokehandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}
	err = cfg.db.UpdateRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Error updating refresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
