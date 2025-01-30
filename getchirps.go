package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getallchirps(w http.ResponseWriter, r *http.Request) {
	dbcps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error getting all chirps: %s", err)
		w.WriteHeader(500)
		return
	}
	var chirps []Chirp
	for _, dbChirp := range dbcps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)
	}

	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshaling chirps: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) getonechirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	cId, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("Error parsing chirp ID: %s", err)
		w.WriteHeader(500)
		return
	}

	dbChirp, err := cfg.db.GetOneChirp(r.Context(), cId)
	if err != nil {
		if err = respondWithError(w, 404, "No chirp found"); err != nil {
			log.Printf("Error responding err: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshaling chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}
