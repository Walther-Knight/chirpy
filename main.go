package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Walther-Knight/chirpy/internal/database"
	"github.com/Walther-Knight/chirpy/internal/middleware"
	"github.com/Walther-Knight/chirpy/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, errDB := sql.Open("postgres", dbURL)
	if errDB != nil {
		log.Printf("Error opening database: %v\n", errDB)
	}
	dbQueries := database.New(db)
	cfg := middleware.ApiConfig{
		Db:          dbQueries,
		Token:       os.Getenv("TOKEN_STRING"),
		PolkaSecret: os.Getenv("POLKA_SECRET"),
	}

	errHttpStart := server.Start(&cfg)
	if errHttpStart != nil {
		log.Printf("Error starting server: %v\n", errHttpStart)
	}

}
