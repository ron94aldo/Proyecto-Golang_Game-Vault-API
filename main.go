package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"game-vault-api/handlers"
)

func main() {
	// Database configuration
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "game_vault"
	}

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

	// Setup routes
	mux := http.NewServeMux()

	// Search endpoints
	mux.HandleFunc("GET /api/search", handlers.SearchGames)
	mux.HandleFunc("GET /api/games/{rawg_id}", handlers.GetGameDetail)

	// Library endpoints
	mux.HandleFunc("GET /api/library", handlers.ListLibrary(db))
	mux.HandleFunc("POST /api/library", handlers.AddGameToLibrary(db))
	mux.HandleFunc("PUT /api/library/{id}", handlers.UpdateGameInLibrary(db))
	mux.HandleFunc("DELETE /api/library/{id}", handlers.DeleteGameFromLibrary(db))

	// Stats endpoint
	mux.HandleFunc("GET /api/library/stats", handlers.GetLibraryStats(db))

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost:8080")

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
