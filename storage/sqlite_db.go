package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDb struct {
	db *sql.DB
}

func New(dbPath string) (*SQLiteDb, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteDb{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT,
		user_id INTEGER NOT NULL,
		is_read BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		notifications_enabled BOOLEAN DEFAULT TRUE
	)`)

	return err
}
