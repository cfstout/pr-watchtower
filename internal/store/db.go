package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(home, ".pr-watchtower.db")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pr_state (
		pr_id INTEGER PRIMARY KEY,
		updated_at TEXT NOT NULL,
		last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

type PRStatus string

const (
	StatusNew     PRStatus = "NEW"
	StatusSeen    PRStatus = "SEEN"
	StatusUpdated PRStatus = "UPDATED"
)

func (s *Store) CheckUpdateStatus(prID int, currentUpdatedAt time.Time) (PRStatus, error) {
	var storedUpdatedAt string
	currentUpdatedAtStr := currentUpdatedAt.Format(time.RFC3339)

	err := s.db.QueryRow("SELECT updated_at FROM pr_state WHERE pr_id = ?", prID).Scan(&storedUpdatedAt)
	if err == sql.ErrNoRows {
		// New PR, insert it
		_, err := s.db.Exec("INSERT INTO pr_state (pr_id, updated_at) VALUES (?, ?)", prID, currentUpdatedAtStr)
		if err != nil {
			return "", err
		}
		return StatusNew, nil
	} else if err != nil {
		// If table schema is old (head_sha), this might fail.
		// For simplicity in this fix, if we error, we might want to try to drop/recreate or just fail.
		// Let's assume we can just fail for now, user can delete db file if needed.
		return "", err
	}

	if storedUpdatedAt != currentUpdatedAtStr {
		// Updated PR, update timestamp
		_, err := s.db.Exec("UPDATE pr_state SET updated_at = ?, last_seen_at = CURRENT_TIMESTAMP WHERE pr_id = ?", currentUpdatedAtStr, prID)
		if err != nil {
			return "", err
		}
		return StatusUpdated, nil
	}

	return StatusSeen, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
