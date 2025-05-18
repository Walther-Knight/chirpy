package main

import (
	"log"

	"github.com/Walther-Knight/chirpy/internal/server"
)

func main() {
	errHttpStart := server.Start()
	if errHttpStart != nil {
		log.Printf("Error starting server: %v\n", errHttpStart)
	}

}
