package repositories

import (
	"database/sql"
	"errors"
	"game-vault-api/models"
	"strconv"
	"time"
)

// PostgresStore implements the LibraryStore interface for PostgreSQL
type PostgresStore struct {
	DB *sql.DB
}

func (s *PostgresStore) ListGames(status string) ([]models.GameLibraryItem, error) {
	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = s.DB.Query("SELECT id, rawg_id, title, genre, platform, cover_url, personal_note, personal_score, status, added_at FROM game_library WHERE status = $1 ORDER BY added_at ASC", status)
	} else {
		rows, err = s.DB.Query("SELECT id, rawg_id, title, genre, platform, cover_url, personal_note, personal_score, status, added_at FROM game_library ORDER BY added_at ASC")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.GameLibraryItem
	for rows.Next() {
		var game models.GameLibraryItem
		var genre, platform, coverURL, personalNote, sStatus sql.NullString
		var personalScore sql.NullInt64

		err := rows.Scan(
			&game.ID,
			&game.RawgID,
			&game.Title,
			&genre,
			&platform,
			&coverURL,
			&personalNote,
			&personalScore,
			&sStatus,
			&game.AddedAt,
		)
		if err != nil {
			return nil, err
		}

		game.Genre = genre.String
		game.Platform = platform.String
		game.CoverURL = coverURL.String
		game.PersonalNote = personalNote.String
		game.PersonalScore = int(personalScore.Int64)
		game.Status = sStatus.String

		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if games == nil {
		games = []models.GameLibraryItem{}
	}

	return games, nil
}

func (s *PostgresStore) AddGame(req models.CreateLibraryRequest) (models.GameLibraryItem, error) {
	var id int
	var addedAt time.Time
	err := s.DB.QueryRow(
		"INSERT INTO game_library (rawg_id, title, genre, platform, cover_url, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, added_at",
		req.RawgID, req.Title, req.Genre, req.Platform, req.CoverURL, "pendiente",
	).Scan(&id, &addedAt)

	if err != nil {
		return models.GameLibraryItem{}, err
	}

	return models.GameLibraryItem{
		ID:       id,
		RawgID:   req.RawgID,
		Title:    req.Title,
		Genre:    req.Genre,
		Platform: req.Platform,
		CoverURL: req.CoverURL,
		Status:   "pendiente",
		AddedAt:  addedAt,
	}, nil
}

// ErrNotFound is returned when a resource is not found in the database
var ErrNotFound = errors.New("resource not found")

func (s *PostgresStore) UpdateGame(id int, req models.UpdateLibraryRequest) error {
	var gameID int
	err := s.DB.QueryRow("SELECT id FROM game_library WHERE id = $1", id).Scan(&gameID)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

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

	if argCount > 1 {
		query = query[:len(query)-2]
	} else {
		return nil
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, id)

	_, err = s.DB.Exec(query, args...)
	return err
}

func (s *PostgresStore) DeleteGame(id int) (bool, error) {
	result, err := s.DB.Exec("DELETE FROM game_library WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (s *PostgresStore) GetStats() (models.LibraryStats, error) {
	var stats models.LibraryStats
	stats.ByStatus = make(map[string]int)

	err := s.DB.QueryRow("SELECT COUNT(*) FROM game_library").Scan(&stats.Total)
	if err != nil {
		return stats, err
	}

	rows, err := s.DB.Query("SELECT status, COUNT(*) FROM game_library GROUP BY status")
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return stats, err
		}
		stats.ByStatus[status] = count
	}

	var avgScore sql.NullFloat64
	err = s.DB.QueryRow("SELECT ROUND(AVG(personal_score), 2) FROM game_library WHERE personal_score IS NOT NULL").Scan(&avgScore)
	if err != nil {
		return stats, err
	}

	if avgScore.Valid {
		stats.AverageScore = avgScore.Float64
	} else {
		stats.AverageScore = 0
	}

	return stats, nil
}
