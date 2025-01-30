package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/NickLiu-0717/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getallchirps(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query().Get("author_id")
	s := r.URL.Query().Get("sort")
	var dbcps []database.Chirp
	var err error
	if a == "" {
		dbcps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			w.WriteHeader(500)
			return
		}
	} else {
		authorID, err := uuid.Parse(s)
		if err != nil {
			log.Printf("Error parsing uuid: %s", err)
			w.WriteHeader(500)
			return
		}
		dbcps, err = cfg.db.GetAuthorChirps(r.Context(), authorID)
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	if s == "desc" {
		sort.Slice(dbcps, func(i, j int) bool { return dbcps[i].CreatedAt.After(dbcps[j].CreatedAt) })
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
	if err = respondWithJSON(w, 200, chirps); err != nil {
		log.Printf("Error responding json: %s", err)
		w.WriteHeader(500)
		return
	}

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
