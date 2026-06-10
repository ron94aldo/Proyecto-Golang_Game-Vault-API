package handlers

import (
	"encoding/json"
	"game-vault-api/models"
	"game-vault-api/repositories"
	"game-vault-api/utils"
	"net/http"
	"strconv"
)

// ListLibrary handles GET /api/library?status={status}
func ListLibrary(store repositories.LibraryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
			return
		}

		status := r.URL.Query().Get("status")
		if status != "" && !utils.ValidateStatus(status) {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid status value")
			return
		}

		games, err := store.ListGames(status)
		if err != nil {
			utils.LogError("Database query", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "500", "Falla de conexión a la BD u otro error interno")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, games)
	}
}

// AddGameToLibrary handles POST /api/library
func AddGameToLibrary(store repositories.LibraryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
			return
		}

		var req models.CreateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid request body")
			return
		}

		if req.RawgID == 0 || req.Title == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Body inválido o campos faltantes")
			return
		}

		game, err := store.AddGame(req)
		if err != nil {
			if err.Error() == "pq: llave duplicada viola restricción de unicidad «game_library_rawg_id_key»" {
				utils.RespondWithError(w, http.StatusConflict, "409", "Recurso duplicado")
				return
			}
			utils.LogError("Database insert", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "500", "Error inesperado del servidor")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, game)
	}
}

// UpdateGameInLibrary handles PUT /api/library/{id}
func UpdateGameInLibrary(store repositories.LibraryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid library ID")
			return
		}

		var req models.UpdateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid request body")
			return
		}

		if req.PersonalScore != nil && !utils.ValidatePersonalScore(*req.PersonalScore) {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Personal score must be between 1 and 10")
			return
		}

		if req.Status != nil && !utils.ValidateStatus(*req.Status) {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid status value")
			return
		}

		err = store.UpdateGame(id, req)
		if err == repositories.ErrNotFound {
			utils.RespondWithError(w, http.StatusNotFound, "404", "Recurso no encontrado en la BD")
			return
		}
		if err != nil {
			// This check captures the 'no fields to update' case without returning an error from DB store
			// But since we want to respond with Bad Request if no fields are provided in req, we should check req first.
			if req.PersonalNote == nil && req.PersonalScore == nil && req.Status == nil {
				utils.RespondWithError(w, http.StatusBadRequest, "400", "No fields to update")
				return
			}

			utils.LogError("Database update", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "500", "Failed to update game")
			return
		}

		// check if no fields were provided explicitly
		if req.PersonalNote == nil && req.PersonalScore == nil && req.Status == nil {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "No fields to update")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Game updated successfully"})
	}
}

// DeleteGameFromLibrary handles DELETE /api/library/{id}
func DeleteGameFromLibrary(store repositories.LibraryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "400", "Invalid library ID")
			return
		}

		deleted, err := store.DeleteGame(id)
		if err != nil {
			utils.LogError("Database delete", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "500", "Failed to delete game")
			return
		}

		if !deleted {
			utils.RespondWithError(w, http.StatusNotFound, "404", "Recurso no encontrado en la BD")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
