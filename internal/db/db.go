package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

var (
	once    sync.Once
	initErr error
)

func InitDB(dbPath, dbName string) error {
	once.Do(func() {
		initErr = initDB(dbPath, dbName)
	})

	return initErr
}

func initDB(dbPath, dbName string) error {
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create database dir: %w", err)
	}

	dbFile := filepath.Join(dbPath, dbName)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	// db.SetMaxOpenConns(30)
	// db.SetConnMaxIdleTime(5)

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA foreign_keys=ON;",
		// "PRAGMA temp_store=MEMORY;",
	}

	for _, pragma := range pragmas {
		_, err := db.Exec(pragma)
		if err != nil {
			_ = db.Close()
			return fmt.Errorf("failed to set pragma '%s': %w", pragma, err)
		}
	}

	tables := []string{
		// Users
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			refresh_token_web TEXT,
  			refresh_token_web_at DATETIME,
  			refresh_token_mobile TEXT,
  			refresh_token_mobile_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// Privates
		`CREATE TABLE IF NOT EXISTS privates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user1_id INTEGER NOT NULL,
			user2_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user1_id, user2_id),
			CHECK(user1_id < user2_id),
			FOREIGN KEY(user1_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(user2_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		// Messages
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_id INTEGER NOT NULL,
			private_id INTEGER,
			message_type TEXT NOT NULL,
			content TEXT NOT NULL,
			delivered INTEGER NOT NULL DEFAULT 0,
			read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(from_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(private_id) REFERENCES privates(id) ON DELETE CASCADE
		);`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			_ = db.Close()
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_messages_private_id ON messages(private_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_from_id ON messages(from_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_privates_user1_id ON privates(user1_id);`,
		`CREATE INDEX IF NOT EXISTS idx_privates_user2_id ON privates(user2_id);`,
	}
	for _, index := range indexes {
		_, err := db.Exec(index)
		if err != nil {
			_ = db.Close()
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	DB = db
	return nil
}

func GetDB() (*sql.DB, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}

	return DB, nil
}

func CloseDB() error {
	if DB == nil {
		return nil
	}

	err := DB.Close()
	if err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}

	DB = nil
	return nil
}
