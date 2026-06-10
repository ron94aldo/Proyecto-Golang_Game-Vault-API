package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"game-vault-api/models"
	"game-vault-api/repositories"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockStore is a mock implementation of repositories.LibraryStore
type MockStore struct {
	ListGamesFunc  func(status string) ([]models.GameLibraryItem, error)
	AddGameFunc    func(req models.CreateLibraryRequest) (models.GameLibraryItem, error)
	UpdateGameFunc func(id int, req models.UpdateLibraryRequest) error
	DeleteGameFunc func(id int) (bool, error)
	GetStatsFunc   func() (models.LibraryStats, error)
}

func (m *MockStore) ListGames(status string) ([]models.GameLibraryItem, error) {
	if m.ListGamesFunc != nil {
		return m.ListGamesFunc(status)
	}
	return nil, nil
}

func (m *MockStore) AddGame(req models.CreateLibraryRequest) (models.GameLibraryItem, error) {
	if m.AddGameFunc != nil {
		return m.AddGameFunc(req)
	}
	return models.GameLibraryItem{}, nil
}

func (m *MockStore) UpdateGame(id int, req models.UpdateLibraryRequest) error {
	if m.UpdateGameFunc != nil {
		return m.UpdateGameFunc(id, req)
	}
	return nil
}

func (m *MockStore) DeleteGame(id int) (bool, error) {
	if m.DeleteGameFunc != nil {
		return m.DeleteGameFunc(id)
	}
	return false, nil
}

func (m *MockStore) GetStats() (models.LibraryStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc()
	}
	return models.LibraryStats{}, nil
}

func TestListLibrary(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		statusParam string
		mockError   error
		wantStatus  int
	}{
		{"Success", "GET", "", nil, http.StatusOK},
		{"Method Not Allowed", "POST", "", nil, http.StatusMethodNotAllowed},
		{"Invalid Status", "GET", "invalid_status", nil, http.StatusBadRequest},
		{"DB Error", "GET", "", errors.New("db error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{
				ListGamesFunc: func(status string) ([]models.GameLibraryItem, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return []models.GameLibraryItem{{ID: 1}}, nil
				},
			}

			url := "/api/library"
			if tt.statusParam != "" {
				url += "?status=" + tt.statusParam
			}
			req, _ := http.NewRequest(tt.method, url, nil)
			rr := httptest.NewRecorder()

			handler := ListLibrary(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}
}

func TestAddGameToLibrary(t *testing.T) {
	validBody := models.CreateLibraryRequest{RawgID: 123, Title: "Test"}

	tests := []struct {
		name       string
		method     string
		body       interface{}
		mockError  error
		wantStatus int
	}{
		{"Success", "POST", validBody, nil, http.StatusCreated},
		{"Method Not Allowed", "GET", validBody, nil, http.StatusMethodNotAllowed},
		{"Invalid JSON", "POST", "{bad json", nil, http.StatusBadRequest},
		{"Missing Fields", "POST", models.CreateLibraryRequest{RawgID: 0}, nil, http.StatusBadRequest},
		{"Conflict Duplicate", "POST", validBody, errors.New("pq: llave duplicada viola restricción de unicidad «game_library_rawg_id_key»"), http.StatusConflict},
		{"DB Error", "POST", validBody, errors.New("db error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{
				AddGameFunc: func(req models.CreateLibraryRequest) (models.GameLibraryItem, error) {
					if tt.mockError != nil {
						return models.GameLibraryItem{}, tt.mockError
					}
					return models.GameLibraryItem{ID: 1}, nil
				},
			}

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(tt.method, "/api/library", bytes.NewBuffer(bodyBytes))
			rr := httptest.NewRecorder()

			handler := AddGameToLibrary(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}
}

func TestUpdateGameInLibrary(t *testing.T) {
	score := 8
	badScore := 15
	statusStr := "completado"
	badStatus := "invalid"

	tests := []struct {
		name       string
		method     string
		id         string
		body       interface{}
		mockError  error
		wantStatus int
	}{
		{"Success", "PUT", "1", models.UpdateLibraryRequest{PersonalScore: &score}, nil, http.StatusOK},
		{"Method Not Allowed", "GET", "1", models.UpdateLibraryRequest{}, nil, http.StatusMethodNotAllowed},
		{"Invalid ID", "PUT", "abc", models.UpdateLibraryRequest{}, nil, http.StatusBadRequest},
		{"Zero ID", "PUT", "0", models.UpdateLibraryRequest{}, nil, http.StatusBadRequest},
		{"Invalid JSON", "PUT", "1", "{bad json", nil, http.StatusBadRequest},
		{"Invalid Score", "PUT", "1", models.UpdateLibraryRequest{PersonalScore: &badScore}, nil, http.StatusBadRequest},
		{"Invalid Status", "PUT", "1", models.UpdateLibraryRequest{Status: &badStatus}, nil, http.StatusBadRequest},
		{"Not Found", "PUT", "1", models.UpdateLibraryRequest{Status: &statusStr}, repositories.ErrNotFound, http.StatusNotFound},
		{"DB Error", "PUT", "1", models.UpdateLibraryRequest{Status: &statusStr}, errors.New("db error"), http.StatusInternalServerError},
		{"No Fields To Update", "PUT", "1", models.UpdateLibraryRequest{}, nil, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{
				UpdateGameFunc: func(id int, req models.UpdateLibraryRequest) error {
					return tt.mockError
				},
			}

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(tt.method, "/api/library/"+tt.id, bytes.NewBuffer(bodyBytes))
			req.SetPathValue("id", tt.id)
			rr := httptest.NewRecorder()

			handler := UpdateGameInLibrary(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}
}

func TestDeleteGameFromLibrary(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		mockReturn bool
		mockError  error
		wantStatus int
	}{
		{"Success", "DELETE", "1", true, nil, http.StatusNoContent},
		{"Method Not Allowed", "GET", "1", false, nil, http.StatusMethodNotAllowed},
		{"Invalid ID", "DELETE", "abc", false, nil, http.StatusBadRequest},
		{"DB Error", "DELETE", "1", false, errors.New("db error"), http.StatusInternalServerError},
		{"Not Found", "DELETE", "1", false, nil, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{
				DeleteGameFunc: func(id int) (bool, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			req, _ := http.NewRequest(tt.method, "/api/library/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			rr := httptest.NewRecorder()

			handler := DeleteGameFromLibrary(mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("%s: got %v want %v", tt.name, status, tt.wantStatus)
			}
		})
	}
}
