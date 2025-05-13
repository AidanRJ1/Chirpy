package main

import (
	"fmt"
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
	mux.HandleFunc("GET /healthz", handlerHealthz)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Serving on Port: %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
