package main

import (
	"log"
	"net/http"

	"github.com/NickLiu-0717/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deletechirp(w http.ResponseWriter, r *http.Request) {
	cID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(cID)
	if err != nil {
		log.Printf("Error parsing chirp ID: %s", err)
		w.WriteHeader(500)
		return
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting access token: %s", err)
		w.WriteHeader(401)
		return
	}
	gotID, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		if err = respondWithError(w, 403, "Forbidden Access"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
		return
	}
	userID, err := cfg.db.GetUserFromChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting user form chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	if userID != gotID {
		if err = respondWithError(w, 403, "Forbidden Access"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
		return
	}
	_, err = cfg.db.DeleteOneChirp(r.Context(), chirpID)
	if err != nil {
		if err = respondWithError(w, 404, "Chirp no found"); err != nil {
			log.Printf("Error responding error: %s", err)
			return
		}
		return
	}
	w.WriteHeader(204)
}
