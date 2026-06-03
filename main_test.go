package main

import (
	"testing"
	"game-vault-api/utils"
)

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
			result := utils.ValidateStatus(tt.status)
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
			result := utils.ValidatePersonalScore(tt.score)
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
			result := utils.ValidateStatus(tt.status)
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
			result := utils.ValidatePersonalScore(tt.score)
			if result != tt.expected {
				t.Errorf("ValidatePersonalScore(%d) = %v, want %v", tt.score, result, tt.expected)
			}
		})
	}
}

// TestGetValidStatuses tests the GetValidStatuses helper function
func TestGetValidStatuses(t *testing.T) {
	statuses := utils.GetValidStatuses()

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
	result1 := utils.ValidatePersonalScore(score)
	result2 := utils.ValidatePersonalScore(score)

	if result1 != result2 {
		t.Errorf("ValidatePersonalScore(%d) returned inconsistent results: %v and %v", score, result1, result2)
	}
}

// TestValidateStatusConsistency tests that ValidateStatus is consistent
func TestValidateStatusConsistency(t *testing.T) {
	// Test multiple calls return same result
	status := "completado"
	result1 := utils.ValidateStatus(status)
	result2 := utils.ValidateStatus(status)

	if result1 != result2 {
		t.Errorf("ValidateStatus(%q) returned inconsistent results: %v and %v", status, result1, result2)
	}
}
