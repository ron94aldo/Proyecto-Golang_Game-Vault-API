package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"game-vault-api/models"
	"game-vault-api/utils"
)

const (
	rawgBaseURL = "https://api.rawg.io/api"
	rawgAPIKey  = "945d345a57cc4c3fb7b4f67211edd4c8"
)

// SearchGames handles GET /api/search?q={nombre}
func SearchGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing_param", "Missing query parameter 'q'")
		return
	}

	// Call RAWG API
	searchURL := rawgBaseURL + "/games?search=" + url.QueryEscape(query) + "&key=" + rawgAPIKey
	
	resp, err := http.Get(searchURL)
	if err != nil {
		utils.LogError("RAWG API call", err)
		utils.RespondWithError(w, http.StatusBadGateway, "external_api_error", "Failed to connect to RAWG API")
		return
	}
	defer resp.Body.Close()

	// Check RAWG response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.LogError("RAWG API error", nil)
		utils.RespondWithError(w, http.StatusBadGateway, "rawg_error", "RAWG API returned an error")
		return
	}

	// Parse RAWG response
	var searchResp models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		utils.LogError("JSON decode", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "parse_error", "Failed to parse RAWG response")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, searchResp)
}

// GetGameDetail handles GET /api/games/{rawg_id}
func GetGameDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Extract rawg_id from URL path
	rawgID := r.PathValue("rawg_id")
	if rawgID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing_param", "Missing rawg_id parameter")
		return
	}

	// Call RAWG API
	gameURL := rawgBaseURL + "/games/" + rawgID + "?key=" + rawgAPIKey
	
	resp, err := http.Get(gameURL)
	if err != nil {
		utils.LogError("RAWG API call", err)
		utils.RespondWithError(w, http.StatusBadGateway, "external_api_error", "Failed to connect to RAWG API")
		return
	}
	defer resp.Body.Close()

	// Check RAWG response status
	if resp.StatusCode == http.StatusNotFound {
		utils.RespondWithError(w, http.StatusNotFound, "game_not_found", "Game not found in RAWG")
		return
	}

	if resp.StatusCode != http.StatusOK {
		utils.LogError("RAWG API error", nil)
		utils.RespondWithError(w, http.StatusBadGateway, "rawg_error", "RAWG API returned an error")
		return
	}

	// Parse RAWG response
	var game models.Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		utils.LogError("JSON decode", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "parse_error", "Failed to parse RAWG response")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, game)
}
