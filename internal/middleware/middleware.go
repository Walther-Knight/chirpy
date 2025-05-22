package middleware

import (
	"bytes"
	"log"
	"net/http"
	"sync/atomic"
	"text/template"

	"github.com/Walther-Knight/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Token          string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

type hitVariables struct {
	HitTotal int32
}

var metricsTemplate = template.Must(template.ParseFiles("./static/templates/admin/metrics.html"))

func (cfg *ApiConfig) HitTotal(w http.ResponseWriter, r *http.Request) {
	hits := hitVariables{
		HitTotal: cfg.FileserverHits.Load(),
	}
	log.Printf("HitTotal endpoint hit. Total Hits: %d\n", hits)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err := metricsTemplate.Execute(&buf, hits)
	if err != nil {
		log.Printf("Error with template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	log.Println("HitTotal Reset endpoint hit.")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	err := cfg.Db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("User table reset failed: %s", err)
	} else {
		log.Println("User table reset.")
	}
}
