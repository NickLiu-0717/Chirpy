package main

import (
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
	secretKey      string
	polkaKey       string
}

type Event struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type userInfo struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type requestBody struct {
	RBody  string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type AccessToken struct {
	Token string `json:"token"`
}
