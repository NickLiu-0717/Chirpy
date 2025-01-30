package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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
