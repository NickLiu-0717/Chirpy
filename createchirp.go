package main

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/NickLiu-0717/Chirpy/internal/auth"
	database "github.com/NickLiu-0717/Chirpy/internal/database"
)

func (cfg *apiConfig) createchirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenSecret, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error geting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}
	gotID, err := auth.ValidateJWT(tokenSecret, cfg.secretKey)
	if err != nil {
		if err = respondWithError(w, 401, "Error: Unauthorized"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	b := requestBody{}
	err = decoder.Decode(&b)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(b.RBody) > 140 {
		if err := respondWithError(w, 400, errorMessage1); err != nil {
			log.Printf("Error responding with error: %s", err)
			w.WriteHeader(500)
			return
		}
		return
	}
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanword := profaneReplace(b.RBody, profaneWords)
	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanword,
		UserID: gotID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	err = respondWithJSON(w, 201, chirp)
	if err != nil {
		log.Printf("Error responding with chirp json: %s", err)
		w.WriteHeader(500)
		return
	}
}
