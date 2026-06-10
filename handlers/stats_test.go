package handlers

import (
	"encoding/json"
	"errors"
	"game-vault-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLibraryStats(t *testing.T) {
	mockStore := &MockStore{
		GetStatsFunc: func() (models.LibraryStats, error) {
			return models.LibraryStats{
				Total: 10,
				ByStatus: map[string]int{
					"completado": 5,
					"jugando":    3,
					"pendiente":  2,
				},
				AverageScore: 8.5,
			}, nil
		},
	}

	req, _ := http.NewRequest("GET", "/api/library/stats", nil)
	rr := httptest.NewRecorder()

	handler := GetLibraryStats(mockStore)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stats models.LibraryStats
	if err := json.NewDecoder(rr.Body).Decode(&stats); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if stats.Total != 10 {
		t.Errorf("expected total 10, got %d", stats.Total)
	}
	if stats.AverageScore != 8.5 {
		t.Errorf("expected avg score 8.5, got %f", stats.AverageScore)
	}
	if stats.ByStatus["completado"] != 5 {
		t.Errorf("expected 5 completado, got %d", stats.ByStatus["completado"])
	}
}

func TestGetLibraryStats_MethodNotAllowed(t *testing.T) {
	mockStore := &MockStore{}
	req, _ := http.NewRequest("POST", "/api/library/stats", nil)
	rr := httptest.NewRecorder()

	handler := GetLibraryStats(mockStore)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("expected %v, got %v", http.StatusMethodNotAllowed, status)
	}
}

func TestGetLibraryStats_StoreError(t *testing.T) {
	mockStore := &MockStore{
		GetStatsFunc: func() (models.LibraryStats, error) {
			return models.LibraryStats{}, errors.New("database error")
		},
	}
	req, _ := http.NewRequest("GET", "/api/library/stats", nil)
	rr := httptest.NewRecorder()

	handler := GetLibraryStats(mockStore)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected %v, got %v", http.StatusInternalServerError, status)
	}
}
