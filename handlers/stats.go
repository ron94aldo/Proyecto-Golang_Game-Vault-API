package handlers

import (
	"database/sql"
	"net/http"
	"game-vault-api/models"
	"game-vault-api/utils"
)

// GetLibraryStats handles GET /api/library/stats
func GetLibraryStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
			return
		}

		var stats models.LibraryStats
		stats.ByStatus = make(map[string]int)

		// Get total games count
		err := db.QueryRow("SELECT COUNT(*) FROM game_library").Scan(&stats.Total)
		if err != nil {
			utils.LogError("Database query total", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to retrieve statistics")
			return
		}

		// Get count by status
		rows, err := db.Query("SELECT status, COUNT(*) FROM game_library GROUP BY status")
		if err != nil {
			utils.LogError("Database query by status", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to retrieve statistics")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var status string
			var count int
			err := rows.Scan(&status, &count)
			if err != nil {
				utils.LogError("Row scan", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to parse statistics")
				return
			}
			stats.ByStatus[status] = count
		}

		// Get average score
		var avgScore sql.NullFloat64
		err = db.QueryRow("SELECT AVG(personal_score) FROM game_library WHERE personal_score IS NOT NULL").Scan(&avgScore)
		if err != nil {
			utils.LogError("Database query average", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to retrieve statistics")
			return
		}

		if avgScore.Valid {
			stats.AverageScore = avgScore.Float64
		} else {
			stats.AverageScore = 0
		}

		utils.RespondWithJSON(w, http.StatusOK, stats)
	}
}
