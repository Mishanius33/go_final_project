package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/nextdate"
)

type TaskEntity struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task TaskEntity
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
			date, err = time.Parse("20060102", task.Date)
			if err != nil {
				sendJSONResponse(w, map[string]string{"error": "Некорректный формат даты"}, http.StatusBadRequest)
				return
			}
		}

		var nextDate string
		if date.Before(time.Now()) {
			if task.Repeat == "" {
				nextDate = time.Now().Format("20060102")
			} else {
				nextDateStr, err := nextdate.NextDate(time.Now(), date.Format("20060102"), task.Repeat)
				if err != nil {
					sendJSONResponse(w, map[string]string{"error": "Некорректный формат повтора: " + task.Repeat}, http.StatusBadRequest)
					return
				}
				nextDate = nextDateStr
			}
		} else {
			nextDate = date.Format("20060102")
		}

		res, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
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
