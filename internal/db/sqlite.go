package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string, reset bool) *sql.DB {

	reset := os.Getenv("APP_ENV") != "production"

	if reset {
		// empty db each start up
		if _, err := os.Stat(path); err == nil {
			err := os.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: improve concurrency on web apps
	_, err = db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		log.Fatal(err)
	}

	schema := `
        CREATE TABLE IF NOT EXISTS lookup (
            id     INTEGER PRIMARY KEY AUTOINCREMENT,
            origin TEXT NOT NULL,
            code   TEXT NOT NULL UNIQUE
        );
        CREATE INDEX IF NOT EXISTS idx_lookup_code ON lookup(code);
    `

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
