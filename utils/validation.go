package utils

// ValidateStatus checks if status is valid
func ValidateStatus(status string) bool {
	validStatuses := map[string]bool{
		"pendiente":   true,
		"jugando":     true,
		"completado":  true,
		"abandonado":  true,
	}
	return validStatuses[status]
}

// ValidatePersonalScore checks if personal score is between 1 and 10
func ValidatePersonalScore(score int) bool {
	return score >= 1 && score <= 10
}

// GetValidStatuses returns the list of valid status values
func GetValidStatuses() []string {
	return []string{"pendiente", "jugando", "completado", "abandonado"}
}
