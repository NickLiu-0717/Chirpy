package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	database "github.com/NickLiu-0717/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
	}
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	apicfg := apiConfig{}
	apicfg.dbQueries = dbQueries
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})
	mux.Handle("/app/", apicfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apicfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validatechirpy)
	server.ListenAndServe()
}
