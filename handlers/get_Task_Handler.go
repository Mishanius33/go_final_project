package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mishanius33/go_final_project/task"

	_ "github.com/mattn/go-sqlite3"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit := 50

		var data interface{}
		var status int
		var err error

		data, status, err = task.GetTasks(db, limit)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
