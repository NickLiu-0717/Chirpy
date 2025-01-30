package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NickLiu-0717/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webhookhandler(w http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	if apikey != cfg.polkaKey {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	defer r.Body.Close()
	var evt Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&evt); err != nil {
		log.Printf("Error decoding json: %s", err)
		w.WriteHeader(500)
		return
	}
	if evt.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	userID, err := uuid.Parse(evt.Data.UserID)
	if err != nil {
		log.Printf("Error parsing uuid string: %s", err)
		w.WriteHeader(500)
		return
	}
	_, err = cfg.db.UpgradeChirpyRed(r.Context(), userID)
	if err != nil {
		if err = respondWithError(w, 404, "User Not Found"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
		return
	}
	w.WriteHeader(204)
}
