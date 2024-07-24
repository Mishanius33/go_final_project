package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createDatabase(dbPath string) error {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
        CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date DATE NOT NULL,
            title TEXT NOT NULL,
            comment TEXT NOT NULL,
            repeat VARCHAR(128) NOT NULL
        )
    `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// Создание индекса по полю date
	_, err = db.Exec("CREATE INDEX idx_scheduler_date ON scheduler (date)")
	if err != nil {
		return err
	}

	return nil
}
