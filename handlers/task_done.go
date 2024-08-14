package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/model"
	"github.com/mishanius33/go_final_project/nextdate"
	"github.com/mishanius33/go_final_project/storage"
)

func TaskDoneHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Нет task ID"})
			return
		}

		resp, _, err := s.GetTaskByID(idStr)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var t model.TaskEntity
		err = json.Unmarshal(resp, &t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if t.Repeat != "" {
			nextDate, err := nextdate.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			err = s.UpdateTaskForDone(model.TaskEntity{
				ID:      t.ID,
				Date:    nextDate,
				Title:   t.Title,
				Comment: t.Comment,
				Repeat:  t.Repeat,
			})
		} else {
			err = s.DeleteTask(t.ID)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct{}{})
	}
}
