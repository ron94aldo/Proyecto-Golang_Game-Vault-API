package models

import "time"

// Game represents a game from RAWG API
type Game struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Released     string `json:"released"`
	BackgroundImage string `json:"background_image"`
	Rating       float64 `json:"rating"`
	Genres       []Genre `json:"genres"`
	Platforms    []Platform `json:"platforms"`
}

// Genre represents a game genre
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Platform represents a game platform
type Platform struct {
	Platform PlatformDetails `json:"platform"`
}

// PlatformDetails contains platform information
type PlatformDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GameLibraryItem represents a game in the user's library
type GameLibraryItem struct {
	ID             int       `json:"id"`
	RawgID         int       `json:"rawg_id"`
	Title          string    `json:"title"`
	Genre          string    `json:"genre"`
	Platform       string    `json:"platform"`
	CoverURL       string    `json:"cover_url"`
	PersonalNote   string    `json:"personal_note"`
	PersonalScore  int       `json:"personal_score"`
	Status         string    `json:"status"`
	AddedAt        time.Time `json:"added_at"`
}

// CreateLibraryRequest represents the request body for adding a game
type CreateLibraryRequest struct {
	RawgID   int    `json:"rawg_id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Platform string `json:"platform"`
	CoverURL string `json:"cover_url"`
}

// UpdateLibraryRequest represents the request body for updating a game
type UpdateLibraryRequest struct {
	PersonalNote  *string `json:"personal_note"`
	PersonalScore *int    `json:"personal_score"`
	Status        *string `json:"status"`
}

// LibraryStats represents statistics of the library
type LibraryStats struct {
	Total        int                `json:"total"`
	ByStatus     map[string]int     `json:"by_status"`
	AverageScore float64            `json:"average_score"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

// SearchResponse represents the RAWG search response
type SearchResponse struct {
	Count   int    `json:"count"`
	Results []Game `json:"results"`
}
