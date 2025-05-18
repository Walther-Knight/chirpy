package main

import (
	"fmt"
	"net/http"
)

func main() {
	newMux := http.NewServeMux()
	httpSrv := &http.Server{
		Handler: newMux,
		Addr:    ":8080",
	}

	err := httpSrv.ListenAndServe()
	if err != nil {
		fmt.Printf("error starting http server: %v", err)
	}

}
