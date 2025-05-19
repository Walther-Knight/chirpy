package middleware

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *ApiConfig) HitTotal(w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	log.Printf("HitTotal endpoint hit. Total Hits: %d\n", hits)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (cfg *ApiConfig) HitReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	log.Println("HitTotal Reset endpoint hit.")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
