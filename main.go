package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	database "github.com/NickLiu-0717/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	apicfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		dev:            platform,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})
	mux.Handle("/app/", apicfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apicfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", apicfg.createchirp)
	mux.HandleFunc("POST /api/users", apicfg.addnewuser)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
