package handlers

import (
	"encoding/json"
	"game-vault-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchGames(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("search") == "error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.Query().Get("search") == "badjson" {
			w.Write([]byte(`{bad json}`))
			return
		}

		searchResp := models.SearchResponse{
			Count: 1,
			Results: []models.Game{
				{ID: 123, Name: "Test Game"},
			},
		}
		json.NewEncoder(w).Encode(searchResp)
	}))
	defer mockServer.Close()

	originalURL := rawgBaseURL
	rawgBaseURL = mockServer.URL
	defer func() { rawgBaseURL = originalURL }()

	tests := []struct {
		name       string
		method     string
		query      string
		wantStatus int
	}{
		{"Method Not Allowed", "POST", "?q=test", http.StatusMethodNotAllowed},
		{"Missing Query Param", "GET", "", http.StatusBadRequest},
		{"Success", "GET", "?q=test", http.StatusOK},
		{"RAWG API Error", "GET", "?q=error", http.StatusBadGateway},
		{"RAWG Bad JSON", "GET", "?q=badjson", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/api/search"+tt.query, nil)
			rr := httptest.NewRecorder()
			SearchGames(rr, req)
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}

	// Test RAWG Connection Error
	rawgBaseURL = "http://invalid-url-that-does-not-exist"
	req, _ := http.NewRequest("GET", "/api/search?q=test", nil)
	rr := httptest.NewRecorder()
	SearchGames(rr, req)
	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("expected BadGateway on connection error, got %v", status)
	}
	rawgBaseURL = mockServer.URL // restore
}

func TestGetGameDetail(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/games/404":
			w.WriteHeader(http.StatusNotFound)
		case "/games/500":
			w.WriteHeader(http.StatusInternalServerError)
		case "/games/badjson":
			w.Write([]byte(`{bad json}`))
		case "/games/123":
			game := models.Game{ID: 123, Name: "Test Game"}
			json.NewEncoder(w).Encode(game)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer mockServer.Close()

	originalURL := rawgBaseURL
	rawgBaseURL = mockServer.URL
	defer func() { rawgBaseURL = originalURL }()

	tests := []struct {
		name       string
		method     string
		rawgID     string
		path       string
		wantStatus int
	}{
		{"Method Not Allowed", "POST", "123", "/api/games/123", http.StatusMethodNotAllowed},
		{"Success", "GET", "123", "/api/games/123", http.StatusOK},
		{"Not Found", "GET", "404", "/api/games/404", http.StatusNotFound},
		{"RAWG Server Error", "GET", "500", "/api/games/500", http.StatusBadGateway},
		{"Bad JSON", "GET", "badjson", "/api/games/badjson", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			if tt.rawgID != "" {
				req.SetPathValue("rawg_id", tt.rawgID)
			}
			rr := httptest.NewRecorder()
			GetGameDetail(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}

	t.Run("Missing ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/games/", nil)
		req.SetPathValue("rawg_id", "")
		rr := httptest.NewRecorder()
		GetGameDetail(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Missing ID: got %v want %v", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("Connection Error", func(t *testing.T) {
		rawgBaseURL = "http://invalid-url-that-does-not-exist"
		req, _ := http.NewRequest("GET", "/api/games/123", nil)
		req.SetPathValue("rawg_id", "123")
		rr := httptest.NewRecorder()
		GetGameDetail(rr, req)
		if status := rr.Code; status != http.StatusBadGateway {
			t.Errorf("Connection error: got %v want %v", status, http.StatusBadGateway)
		}
	})
}
