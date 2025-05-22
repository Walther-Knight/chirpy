package server

import (
	"log"
	"net/http"

	"github.com/Walther-Knight/chirpy/internal/api"
	"github.com/Walther-Knight/chirpy/internal/middleware"
)

func Start(cfg *middleware.ApiConfig) error {
	newMux := http.NewServeMux()
	httpSrv := &http.Server{
		Handler: newMux,
		Addr:    ":8080",
	}
	log.Println("Starting handlers...")
	//admin functions
	newMux.HandleFunc("GET /api/healthz", api.Health)
	newMux.HandleFunc("GET /admin/metrics", cfg.HitTotal)
	newMux.HandleFunc("POST /admin/reset", cfg.Reset)
	//application functions
	newMux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) { api.UserLogin(cfg, w, r) })
	newMux.HandleFunc("GET /api/chirps/{id}", func(w http.ResponseWriter, r *http.Request) { api.GetChirp(cfg, w, r) })
	newMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) { api.NewChirp(cfg, w, r) })
	newMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) { api.GetAllChirps(cfg, w, r) })
	newMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) { api.NewUser(cfg, w, r) })
	newMux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}
