package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"game-vault-api/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	rr := httptest.NewRecorder()
	RespondWithError(rr, http.StatusBadRequest, "400", "Bad Request")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var response models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Code != "400" {
		t.Errorf("Expected code '400', got %s", response.Code)
	}
	if response.Error != "Bad Request" {
		t.Errorf("Expected error 'Bad Request', got %s", response.Error)
	}
}

func TestRespondWithJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	RespondWithJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got %s", response["message"])
	}
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) // Reset after test

	err := errors.New("test error")
	LogError("TestContext", err)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "[ERROR] TestContext: test error") {
		t.Errorf("Expected log output to contain '[ERROR] TestContext: test error', got %s", logOutput)
	}

	// Test with nil error (should not log)
	buf.Reset()
	LogError("TestContextNil", nil)
	if buf.Len() > 0 {
		t.Errorf("Expected no log output for nil error, got %s", buf.String())
	}
}

// TestValidateStatus tests the ValidateStatus function
func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		// Valid statuses
		{"Valid status: pendiente", "pendiente", true},
		{"Valid status: jugando", "jugando", true},
		{"Valid status: completado", "completado", true},
		{"Valid status: abandonado", "abandonado", true},

		// Invalid statuses
		{"Invalid status: completed", "completed", false},
		{"Invalid status: playing", "playing", false},
		{"Invalid status: empty", "", false},
		{"Invalid status: unknown", "unknown", false},
		{"Invalid status: PENDIENTE uppercase", "PENDIENTE", false},
		{"Invalid status: pending english", "pending", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidateStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

// TestValidatePersonalScore tests the ValidatePersonalScore function
func TestValidatePersonalScore(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected bool
	}{
		// Valid scores
		{"Valid score: 1", 1, true},
		{"Valid score: 5", 5, true},
		{"Valid score: 10", 10, true},
		{"Valid score: 7", 7, true},
		{"Valid score: 3", 3, true},

		// Invalid scores - below range
		{"Invalid score: 0", 0, false},
		{"Invalid score: -1", -1, false},
		{"Invalid score: -10", -10, false},

		// Invalid scores - above range
		{"Invalid score: 11", 11, false},
		{"Invalid score: 100", 100, false},
		{"Invalid score: 15", 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePersonalScore(tt.score)
			if result != tt.expected {
				t.Errorf("ValidatePersonalScore(%d) = %v, want %v", tt.score, result, tt.expected)
			}
		})
	}
}

// Additional tests for edge cases and coverage

// TestValidateStatusEdgeCases tests edge cases for status validation
func TestValidateStatusEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Empty string", "", false},
		{"Whitespace", " ", false},
		{"Multiple spaces", "  pendiente  ", false},
		{"Case sensitivity", "Pendiente", false},
		{"Partial match", "pend", false},
		{"Typo", "pndiente", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatus(tt.status)
			if result != tt.expected {
				t.Errorf("ValidateStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

// TestValidatePersonalScoreBoundaries tests boundary conditions
func TestValidatePersonalScoreBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected bool
	}{
		{"Lower boundary: 1", 1, true},
		{"Upper boundary: 10", 10, true},
		{"Just below lower: 0", 0, false},
		{"Just above upper: 11", 11, false},
		{"Far below: -100", -100, false},
		{"Far above: 1000", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePersonalScore(tt.score)
			if result != tt.expected {
				t.Errorf("ValidatePersonalScore(%d) = %v, want %v", tt.score, result, tt.expected)
			}
		})
	}
}

// TestGetValidStatuses tests the GetValidStatuses helper function
func TestGetValidStatuses(t *testing.T) {
	statuses := GetValidStatuses()

	expectedCount := 4
	if len(statuses) != expectedCount {
		t.Errorf("GetValidStatuses() returned %d statuses, want %d", len(statuses), expectedCount)
	}

	expectedStatuses := map[string]bool{
		"pendiente":  true,
		"jugando":    true,
		"completado": true,
		"abandonado": true,
	}

	for _, status := range statuses {
		if !expectedStatuses[status] {
			t.Errorf("GetValidStatuses() returned unexpected status: %s", status)
		}
	}
}

// TestValidatePersonalScoreConsistency tests that ValidatePersonalScore is consistent
func TestValidatePersonalScoreConsistency(t *testing.T) {
	// Test multiple calls return same result
	score := 7
	result1 := ValidatePersonalScore(score)
	result2 := ValidatePersonalScore(score)

	if result1 != result2 {
		t.Errorf("ValidatePersonalScore(%d) returned inconsistent results: %v and %v", score, result1, result2)
	}
}

// TestValidateStatusConsistency tests that ValidateStatus is consistent
func TestValidateStatusConsistency(t *testing.T) {
	// Test multiple calls return same result
	status := "completado"
	result1 := ValidateStatus(status)
	result2 := ValidateStatus(status)

	if result1 != result2 {
		t.Errorf("ValidateStatus(%q) returned inconsistent results: %v and %v", status, result1, result2)
	}
}
