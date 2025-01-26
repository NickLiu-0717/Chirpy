package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	database "github.com/NickLiu-0717/Chirpy/internal/database"
	"github.com/google/uuid"
)

const errorMessage1 = "Chirp is too long"

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	dev            string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func profaneReplace(str string, profaneWords []string) string {
	lowerstr := strings.ToLower(str)
	splitOrigin := strings.Split(str, " ")
	for _, word := range profaneWords {
		lowerstr = strings.ReplaceAll(lowerstr, word, "****")
	}
	splitLower := strings.Split(lowerstr, " ")
	for idx, word := range splitLower {
		if word == "****" {
			splitOrigin[idx] = "****"
		}
	}
	return strings.Join(splitOrigin, " ")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlTemplate := `
    <html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    </html>
    `
	htmlResponse := fmt.Sprintf(htmlTemplate, cfg.fileserverHits.Load())
	w.Write([]byte(htmlResponse))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.dev != "dev" {
		w.WriteHeader(403)
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting all users: %s", err)
		w.WriteHeader(500)
		return
	}
}

func (cfg *apiConfig) addnewuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type newUser struct {
		Email string `json:"email"`
	}
	defer r.Body.Close()
	params := newUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) createchirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type requestBody struct {
		RBody  string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type Chirps struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	b := requestBody{}
	err := decoder.Decode(&b)
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
		UserID: b.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	chirp := Chirps{
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
