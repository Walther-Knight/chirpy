package server

import (
	"log"
	"net/http"
)

func Start() error {
	newMux := http.NewServeMux()
	httpSrv := &http.Server{
		Handler: newMux,
		Addr:    ":8080",
	}
	log.Println("Starting handlers...")
	newMux.HandleFunc("/healthz", health)
	newMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
