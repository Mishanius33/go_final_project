package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type TaskEntity struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit := 50

		var data interface{}
		var status int
		var err error

		data, status, err = getTasks(db, limit)

		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func getTasks(db *sql.DB, limit int) (interface{}, int, error) {
	var tasks []TaskEntity

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var t TaskEntity
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(tasks) == 0 {
		tasks = []TaskEntity{}
	}

	return map[string][]TaskEntity{"tasks": tasks}, http.StatusOK, nil
}
