package handlers

import (
	"database/sql"
	"net/http"

	"github.com/mishanius33/go_final_project/task"
)

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				data, status, err := task.GetTaskByID(db, id)
				if err != nil {
					http.Error(w, string(data), status)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(status)
				w.Write(data)
			} else {
				http.Error(w, `{"error": "Идентификатор не указан"}`, http.StatusBadRequest)
			}
		case http.MethodPost:
			AddTaskHandler(db)(w, r)
		case http.MethodPut:
			EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r)
		}
	}
}
