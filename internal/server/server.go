package server

import (
	"log"
	"net/http"

	"github.com/Walther-Knight/chirpy/internal/middleware"
)

func Start() error {
	newMux := http.NewServeMux()
	httpSrv := &http.Server{
		Handler: newMux,
		Addr:    ":8080",
	}
	var hitCount middleware.ApiConfig
	log.Println("Starting handlers...")
	newMux.HandleFunc("GET /healthz", health)
	newMux.HandleFunc("GET /metrics", hitCount.HitTotal)
	newMux.HandleFunc("POST /reset", hitCount.HitReset)
	newMux.Handle("/app/", hitCount.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
