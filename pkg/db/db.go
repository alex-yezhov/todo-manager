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

var DB *sql.DB

func getDBPath() string {
	path := strings.TrimSpace(os.Getenv(envDBFile))
	if path == "" {
		return defaultDBFile
	}
	return path
}

func InitDB() (*sql.DB, error) {
	dbPath := getDBPath()

	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать папку для базы: %w", err)
		}
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу: %w", err)
	}

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("не удалось подключиться к базе: %w", err)
	}

	if err := createTable(database); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	DB = database
	return database, nil
}

func createTable(database *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT NOT NULL DEFAULT '',
		repeat TEXT NOT NULL DEFAULT ''
	);`

	_, err := database.Exec(query)
	return err
}
