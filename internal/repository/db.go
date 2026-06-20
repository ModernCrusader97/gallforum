package repository

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(migration)); err != nil {
		return nil, err
	}
	return db, nil
}
