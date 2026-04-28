package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	defaultDBFile = "scheduler.db"
	envDBFile     = "TODO_DBFILE"
)

type Store struct {
	db *sql.DB
}

func getDBPath() string {
	path := strings.TrimSpace(os.Getenv(envDBFile))
	if path == "" {
		return defaultDBFile
	}
	return path
}

func InitDB() (*Store, error) {
	dbPath := getDBPath()

	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("не удалось создать папку для базы: %w", err)
		}
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("не удалось подключиться к базе: %w", err)
	}

	if err := createTable(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	return &Store{db: conn}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func createTable(conn *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL CHECK(length(date) = 8),
		title TEXT NOT NULL,
		comment TEXT NOT NULL DEFAULT '',
		repeat TEXT NOT NULL DEFAULT '' CHECK(length(repeat) <= 128)
	);

	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
	`

	_, err := conn.Exec(query)
	return err
}
