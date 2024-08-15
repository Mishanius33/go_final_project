package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/model"
	"github.com/mishanius33/go_final_project/nextdate"
	"github.com/mishanius33/go_final_project/storage"
)

func AddTaskHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task model.TaskRequest
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Ошибка JSON: "+err.Error())
			return
		}

		if task.Title == "" {
			respondWithError(w, http.StatusBadRequest, "Требуется title")
			return
		}

		var date time.Time
		if task.Date == "" || task.Date == time.Now().Format(DateFormat) {
			date = time.Now()
		} else {
			date, err = time.Parse(DateFormat, task.Date)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Неверный формат даты")
				return
			}
		}

		var nextDate string
		if date.Before(time.Now()) {
			if task.Repeat == "" {
				nextDate = time.Now().Format(DateFormat)
			} else {
				nextDateStr, err := nextdate.NextDate(time.Now(), date.Format(DateFormat), task.Repeat)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Неверный формат повторения: "+task.Repeat)
					return
				}
				nextDate = nextDateStr
			}
		} else {
			nextDate = date.Format(DateFormat)
		}

		id, err := s.InsertTask(nextDate, task.Title, task.Comment, task.Repeat)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка добавления задачи: "+err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, model.TaskResponse{ID: id})
	}
}
