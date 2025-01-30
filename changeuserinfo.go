package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NickLiu-0717/Chirpy/internal/auth"
	"github.com/NickLiu-0717/Chirpy/internal/database"
)

func (cfg *apiConfig) changeUserInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting access token: %s", err)
		w.WriteHeader(500)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		if err = respondWithError(w, 401, "Error: Unauthorized"); err != nil {
			log.Printf("Error responding error: %s", err)
			w.WriteHeader(500)
			return
		}
	}
	uInfo := userInfo{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&uInfo); err != nil {
		log.Printf("Error decoding json: %s", err)
		w.WriteHeader(500)
		return
	}
	hashed_password, err := auth.HashPassword(uInfo.Password)
	if err != nil {
		log.Printf("Error making hased password: %s", err)
		w.WriteHeader(500)
		return
	}
	cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: hashed_password,
		Email:          uInfo.Email,
		ID:             userID,
	})
	dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user by ID: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	if err = respondWithJSON(w, 200, user); err != nil {
		log.Printf("Error responding json: %s", err)
		w.WriteHeader(500)
		return
	}
}
