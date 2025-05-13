package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
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
