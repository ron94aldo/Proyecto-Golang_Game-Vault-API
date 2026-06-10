package main

import (
	"database/sql"
	"fmt"
	"game-vault-api/handlers"
	"game-vault-api/repositories"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func main() {
	// Database configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "ronaldo")
	dbName := getEnv("DB_NAME", "game_vault")

	// Connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Create table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS game_library (
			id             SERIAL PRIMARY KEY,
			rawg_id        INTEGER NOT NULL UNIQUE,
			title          VARCHAR(255) NOT NULL,
			genre          VARCHAR(100),
			platform       VARCHAR(100),
			cover_url      TEXT,
			personal_note  TEXT,
			personal_score INTEGER,
			status         VARCHAR(20),
			added_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_rawg_id ON game_library(rawg_id);
		CREATE INDEX IF NOT EXISTS idx_status ON game_library(status);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Database table verified/created successfully")

	store := &repositories.PostgresStore{DB: db}

	// Setup routes
	mux := http.NewServeMux()

	// Search endpoints
	mux.HandleFunc("GET /api/search", handlers.SearchGames)
	mux.HandleFunc("GET /api/games/{rawg_id}", handlers.GetGameDetail)

	// Library endpoints
	mux.HandleFunc("GET /api/library", handlers.ListLibrary(store))
	mux.HandleFunc("POST /api/library", handlers.AddGameToLibrary(store))
	mux.HandleFunc("PUT /api/library/{id}", handlers.UpdateGameInLibrary(store))
	mux.HandleFunc("DELETE /api/library/{id}", handlers.DeleteGameFromLibrary(store))

	// Stats endpoint
	mux.HandleFunc("GET /api/library/stats", handlers.GetLibraryStats(store))

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost:8080")

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
