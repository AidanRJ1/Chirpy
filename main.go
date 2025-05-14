package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AidanRJ1/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error Reading Database: %v", err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", handler)))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Serving on Port: %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
