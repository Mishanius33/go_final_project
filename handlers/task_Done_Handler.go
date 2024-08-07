package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/nextdate"
	"github.com/mishanius33/go_final_project/task"
)

var tasks []TaskEntity

func TaskDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Нет task ID"})
			return
		}

		resp, _, err := task.GetTaskByID(db, idStr)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var t TaskEntity
		err = json.Unmarshal(resp, &t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		t.Done = true
		if t.Repeat != "" {
			nextDate, err := nextdate.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			t.Date = nextDate
			err = UpdateTask(db, t)
		} else {
			err = DeleteTask(db, t.ID)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func UpdateTask(db *sql.DB, task TaskEntity) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?, done = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat, task.Done, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(db *sql.DB, taskID string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskID)
	if err != nil {
		return err
	}

	return nil
}
