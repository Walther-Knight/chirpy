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
	newMux.Handle("/", http.FileServer(http.Dir("./static")))

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	return httpSrv.ListenAndServe()
}
