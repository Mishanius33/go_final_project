package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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
            comment TEXT,
            repeat VARCHAR(128)
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

func updateTaskDate(db *sql.DB, taskID int) error {
	var repeatRule string
	var repeatDate string

	err := db.QueryRow("SELECT repeat, date FROM scheduler WHERE id = ?", taskID).Scan(&repeatRule, &repeatDate)
	if err != nil {
		return err
	}

	// Парсим дату задачи
	date, err := strconv.ParseInt(repeatDate, 10, 64)
	if err != nil {
		return err
	}
	taskTime := time.Unix(date, 0)

	var newTaskTime time.Time
	switch {
	case strings.HasPrefix(repeatRule, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeatRule, "d "))
		if err != nil || days < 1 || days > 400 {
			return fmt.Errorf("Неверное правило повторения: %s", repeatRule)
		}
		newTaskTime = taskTime.AddDate(0, 0, days)
	case repeatRule == "y":
		newTaskTime = taskTime.AddDate(1, 0, 0)
	default:
		return fmt.Errorf("Неверное правило повторения: %s", repeatRule)
	}

	if repeatRule == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", taskID)
		return err
	}

	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newTaskTime.Unix(), taskID)
	return err
}
