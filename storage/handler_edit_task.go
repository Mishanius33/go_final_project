package storage

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/common"
	"github.com/mishanius33/go_final_project/nextdate"
)

func EditTaskHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task common.TaskEntity
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": "Ошибка при декодировании JSON: " + err.Error()}, http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			sendJSONResponse(w, map[string]string{"error": "ID задачи не указан"}, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			sendJSONResponse(w, map[string]string{"error": "Заголовок задачи не указан"}, http.StatusBadRequest)
			return
		}

		var date time.Time
		if task.Date == "" {
			date = time.Now()
		} else {
			date, err = time.Parse(common.DateFormat, task.Date)
			if err != nil {
				sendJSONResponse(w, map[string]string{"error": "Некорректный формат даты"}, http.StatusBadRequest)
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
					sendJSONResponse(w, map[string]string{"error": "Некорректный формат повтора: " + task.Repeat}, http.StatusBadRequest)
					return
				}
				nextDate = nextDateStr
			}
		} else {
			nextDate = date.Format(common.DateFormat)
		}

		res, err := s.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			nextDate, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": "Ошибка при обновлении задачи: " + err.Error()}, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": "Ошибка при получении количества затронутых строк: " + err.Error()}, http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			sendJSONResponse(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
			return
		}

		sendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
	}
}
