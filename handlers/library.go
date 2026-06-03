package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"game-vault-api/models"
	"game-vault-api/utils"
)

// ListLibrary handles GET /api/library?status={status}
func ListLibrary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
			return
		}

		status := r.URL.Query().Get("status")
		var rows *sql.Rows
		var err error

		if status != "" {
			// Validate status if provided
			if !utils.ValidateStatus(status) {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid_status", "Invalid status value")
				return
			}
			rows, err = db.Query("SELECT id, rawg_id, title, genre, platform, cover_url, personal_note, personal_score, status, added_at FROM game_library WHERE status = $1 ORDER BY added_at DESC", status)
		} else {
			rows, err = db.Query("SELECT id, rawg_id, title, genre, platform, cover_url, personal_note, personal_score, status, added_at FROM game_library ORDER BY added_at DESC")
		}

		if err != nil {
			utils.LogError("Database query", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to retrieve library")
			return
		}
		defer rows.Close()

		var games []models.GameLibraryItem
		for rows.Next() {
			var game models.GameLibraryItem
			err := rows.Scan(&game.ID, &game.RawgID, &game.Title, &game.Genre, &game.Platform, &game.CoverURL, &game.PersonalNote, &game.PersonalScore, &game.Status, &game.AddedAt)
			if err != nil {
				utils.LogError("Row scan", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to parse library data")
				return
			}
			games = append(games, game)
		}

		if err = rows.Err(); err != nil {
			utils.LogError("Rows error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to retrieve library")
			return
		}

		if games == nil {
			games = []models.GameLibraryItem{}
		}

		utils.RespondWithJSON(w, http.StatusOK, games)
	}
}

// AddGameToLibrary handles POST /api/library
func AddGameToLibrary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
			return
		}

		var req models.CreateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
			return
		}

		// Validate required fields
		if req.RawgID == 0 || req.Title == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing_fields", "Missing required fields: rawg_id, title")
			return
		}

		// Insert into database
		var id int
		err := db.QueryRow(
			"INSERT INTO game_library (rawg_id, title, genre, platform, cover_url, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			req.RawgID, req.Title, req.Genre, req.Platform, req.CoverURL, "pendiente",
		).Scan(&id)

		if err != nil {
			// Check for duplicate entry
			if err.Error() == "pq: duplicate key value violates unique constraint \"game_library_rawg_id_key\"" {
				utils.RespondWithError(w, http.StatusConflict, "duplicate_game", "Game already exists in library")
				return
			}
			utils.LogError("Database insert", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to add game to library")
			return
		}

		game := models.GameLibraryItem{
			ID:     id,
			RawgID: req.RawgID,
			Title:  req.Title,
			Genre:  req.Genre,
			Platform: req.Platform,
			CoverURL: req.CoverURL,
			Status: "pendiente",
		}

		utils.RespondWithJSON(w, http.StatusCreated, game)
	}
}

// UpdateGameInLibrary handles PUT /api/library/{id}
func UpdateGameInLibrary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid library ID")
			return
		}

		var req models.UpdateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
			return
		}

		// Check if game exists
		var gameID int
		err = db.QueryRow("SELECT id FROM game_library WHERE id = $1", id).Scan(&gameID)
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "game_not_found", "Game not found in library")
			return
		}
		if err != nil {
			utils.LogError("Database query", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to update game")
			return
		}

		// Validate fields if provided
		if req.PersonalScore != nil && !utils.ValidatePersonalScore(*req.PersonalScore) {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_score", "Personal score must be between 1 and 10")
			return
		}

		if req.Status != nil && !utils.ValidateStatus(*req.Status) {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_status", "Invalid status value")
			return
		}

		// Build dynamic update query
		var args []interface{}
		query := "UPDATE game_library SET "
		argCount := 1

		if req.PersonalNote != nil {
			query += "personal_note = $" + strconv.Itoa(argCount) + ", "
			args = append(args, *req.PersonalNote)
			argCount++
		}

		if req.PersonalScore != nil {
			query += "personal_score = $" + strconv.Itoa(argCount) + ", "
			args = append(args, *req.PersonalScore)
			argCount++
		}

		if req.Status != nil {
			query += "status = $" + strconv.Itoa(argCount) + ", "
			args = append(args, *req.Status)
			argCount++
		}

		// Remove trailing comma and space
		if argCount > 1 {
			query = query[:len(query)-2]
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "no_fields", "No fields to update")
			return
		}

		query += " WHERE id = $" + strconv.Itoa(argCount)
		args = append(args, id)

		_, err = db.Exec(query, args...)
		if err != nil {
			utils.LogError("Database update", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to update game")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Game updated successfully"})
	}
}

// DeleteGameFromLibrary handles DELETE /api/library/{id}
func DeleteGameFromLibrary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid library ID")
			return
		}

		result, err := db.Exec("DELETE FROM game_library WHERE id = $1", id)
		if err != nil {
			utils.LogError("Database delete", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to delete game")
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			utils.LogError("Rows affected", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "db_error", "Failed to verify deletion")
			return
		}

		if rowsAffected == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "game_not_found", "Game not found in library")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
