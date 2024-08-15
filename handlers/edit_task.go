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
			respondWithError(w, http.StatusBadRequest, "Ошибка при декодировании JSON: "+err.Error())
			return
		}

		if task.ID == "" {
			respondWithError(w, http.StatusBadRequest, "ID задачи не указан")
			return
		}

		if task.Title == "" {
			respondWithError(w, http.StatusBadRequest, "Заголовок задачи не указан")

			return
		}

		var date time.Time
		if task.Date == "" {
			date = time.Now()
		} else {
			date, err = time.Parse(DateFormat, task.Date)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Некорректный формат даты")
				return
			}
		}

		nextDateStr, err := nextdate.NextDate(time.Now(), date.Format(DateFormat), task.Repeat)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректный формат повтора: "+task.Repeat)
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
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи: "+err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]interface{}{})
	}
}
