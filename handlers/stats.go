package handlers

import (
	"game-vault-api/repositories"
	"game-vault-api/utils"
	"net/http"
)

// GetLibraryStats handles GET /api/library/stats
func GetLibraryStats(store repositories.LibraryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
			return
		}

		stats, err := store.GetStats()
		if err != nil {
			utils.LogError("Database query stats", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "500", "Failed to retrieve statistics")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, stats)
	}
}
