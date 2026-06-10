package repositories

import "game-vault-api/models"

// LibraryStore defines the interface for database operations
type LibraryStore interface {
	ListGames(status string) ([]models.GameLibraryItem, error)
	AddGame(req models.CreateLibraryRequest) (models.GameLibraryItem, error)
	UpdateGame(id int, req models.UpdateLibraryRequest) error
	DeleteGame(id int) (bool, error)
	GetStats() (models.LibraryStats, error)
}
