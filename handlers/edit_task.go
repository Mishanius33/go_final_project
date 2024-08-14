package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/model"
	"github.com/mishanius33/go_final_project/nextdate"
	"github.com/mishanius33/go_final_project/storage"
)

func EditTaskHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task model.TaskEntity
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
			date, err = time.Parse(DateFormat, task.Date)
			if err != nil {
				sendJSONResponse(w, map[string]string{"error": "Некорректный формат даты"}, http.StatusBadRequest)
				return
			}
		}

		nextDateStr, err := nextdate.NextDate(time.Now(), date.Format(DateFormat), task.Repeat)
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": "Некорректный формат повтора: " + task.Repeat}, http.StatusBadRequest)
			return
		}

		err = s.UpdateTask(model.TaskEntity{
			ID:      task.ID,
			Title:   task.Title,
			Comment: task.Comment,
			Date:    nextDateStr,
			Repeat:  task.Repeat,
		})
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": "Ошибка при обновлении задачи: " + err.Error()}, http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
	}
}
