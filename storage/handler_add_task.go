package storage

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/common"
	"github.com/mishanius33/go_final_project/nextdate"
)

func AddTaskHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task common.TaskRequest
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			sendJSONResponse(w, common.TaskResponse{Err: "Ошибка JSON: " + err.Error()}, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			sendJSONResponse(w, common.TaskResponse{Err: "Требуется title"}, http.StatusBadRequest)
			return
		}

		var date time.Time
		if task.Date == "" || task.Date == time.Now().Format(common.DateFormat) {
			date = time.Now()
		} else {
			date, err = time.Parse(common.DateFormat, task.Date)
			if err != nil {
				sendJSONResponse(w, common.TaskResponse{Err: "Неверный формат даты"}, http.StatusBadRequest)
				return
			}
		}

		var nextDate string
		if date.Before(time.Now()) {
			if task.Repeat == "" {
				nextDate = time.Now().Format(common.DateFormat)
			} else {
				nextDateStr, err := nextdate.NextDate(time.Now(), date.Format(common.DateFormat), task.Repeat)
				if err != nil {
					sendJSONResponse(w, common.TaskResponse{Err: "Неверный формат повторения: " + task.Repeat}, http.StatusBadRequest)
					return
				}
				nextDate = nextDateStr
			}
		} else {
			nextDate = date.Format(common.DateFormat)
		}

		res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", nextDate, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, common.TaskResponse{ID: id}, http.StatusOK)
	}
}
