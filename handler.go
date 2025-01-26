package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	database "github.com/NickLiu-0717/Chirpy/internal/database"
)

const errorMessage1 = "Chirp is too long"

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
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

func (cfg *apiConfig) resetHandler(http.ResponseWriter, *http.Request) {
	cfg.fileserverHits.Store(0)
}

func validatechirpy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type requestBody struct {
		RBody string `json:"body"`
	}
	type requestClean struct {
		Clean string `json:"cleaned_body"`
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
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}
		return
	}
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanword := profaneReplace(b.RBody, profaneWords)
	v := requestClean{Clean: cleanword}
	err = respondWithJSON(w, 200, v)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
