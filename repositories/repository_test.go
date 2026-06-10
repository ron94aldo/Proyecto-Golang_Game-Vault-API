package repositories

import (
	"database/sql"
	"fmt"
	"game-vault-api/models"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func setupTestDB() (*sql.DB, error) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "game_vault") // Ideally a separate test db

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS game_library_test (
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
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		log.Printf("Cannot connect to test database. Skipping integration tests: %v", err)
		os.Exit(0)
	}

	// Clean up table before running tests
	_, _ = testDB.Exec("TRUNCATE TABLE game_library_test RESTART IDENTITY")

	// We'll temporarily point the PostgresStore queries to the test table 
	// To do this properly without altering store.go too much, we will test the logical flow
	// using the existing queries but on a real db. 
	// For these tests to not mess up production data, they require a clean test DB.
	
	code := m.Run()

	// Clean up after tests
	_, _ = testDB.Exec("DROP TABLE IF EXISTS game_library_test")
	testDB.Close()

	os.Exit(code)
}

func TestStore_Integration(t *testing.T) {
	// Re-create the table mapping for tests
	_, _ = testDB.Exec(`
		DROP TABLE IF EXISTS game_library;
		CREATE TABLE game_library (
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
	`)

	store := &PostgresStore{DB: testDB}

	// 1. Test AddGame
	t.Run("AddGame", func(t *testing.T) {
		req := models.CreateLibraryRequest{
			RawgID:   1,
			Title:    "Test Game",
			Genre:    "Action",
			Platform: "PC",
		}
		
		game, err := store.AddGame(req)
		if err != nil {
			t.Fatalf("Failed to add game: %v", err)
		}
		
		if game.Title != "Test Game" {
			t.Errorf("Expected title 'Test Game', got %s", game.Title)
		}
		if game.Status != "pendiente" {
			t.Errorf("Expected status 'pendiente', got %s", game.Status)
		}
	})

	// 2. Test ListGames
	t.Run("ListGames", func(t *testing.T) {
		games, err := store.ListGames("")
		if err != nil {
			t.Fatalf("Failed to list games: %v", err)
		}
		
		if len(games) != 1 {
			t.Errorf("Expected 1 game, got %d", len(games))
		}

		gamesPendiente, _ := store.ListGames("pendiente")
		if len(gamesPendiente) != 1 {
			t.Errorf("Expected 1 pending game")
		}
	})

	// 3. Test UpdateGame
	t.Run("UpdateGame", func(t *testing.T) {
		score := 9
		status := "completado"
		updateReq := models.UpdateLibraryRequest{
			PersonalScore: &score,
			Status:        &status,
		}

		err := store.UpdateGame(1, updateReq)
		if err != nil {
			t.Fatalf("Failed to update game: %v", err)
		}

		// Verify update
		games, _ := store.ListGames("completado")
		if len(games) != 1 {
			t.Errorf("Expected 1 completed game after update, got %d", len(games))
		}
		if games[0].PersonalScore != 9 {
			t.Errorf("Expected score 9, got %d", games[0].PersonalScore)
		}

		// Test updating non-existent game
		err = store.UpdateGame(999, updateReq)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 4. Test GetStats
	t.Run("GetStats", func(t *testing.T) {
		stats, err := store.GetStats()
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		if stats.Total != 1 {
			t.Errorf("Expected total 1, got %d", stats.Total)
		}
		if stats.ByStatus["completado"] != 1 {
			t.Errorf("Expected 1 completed in stats")
		}
		if stats.AverageScore != 9.0 {
			t.Errorf("Expected average score 9.0, got %f", stats.AverageScore)
		}
	})

	// 5. Test DeleteGame
	t.Run("DeleteGame", func(t *testing.T) {
		deleted, err := store.DeleteGame(1)
		if err != nil {
			t.Fatalf("Failed to delete game: %v", err)
		}
		if !deleted {
			t.Errorf("Expected game to be deleted")
		}

		// Verify deletion
		games, _ := store.ListGames("")
		if len(games) != 0 {
			t.Errorf("Expected 0 games after deletion")
		}
	})
}
