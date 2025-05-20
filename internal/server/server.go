package server

import (
	"log"
	"net/http"

	"github.com/Walther-Knight/chirpy/internal/api"
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
	//admin functions
	newMux.HandleFunc("GET /api/healthz", api.Health)
	newMux.HandleFunc("GET /admin/metrics", hitCount.HitTotal)
	newMux.HandleFunc("POST /admin/reset", hitCount.HitReset)
	//application functions
	newMux.Handle("/app/", hitCount.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))
	newMux.HandleFunc("POST /api/validate_chirp", api.Validate)

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}
