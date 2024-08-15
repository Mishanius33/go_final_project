package handlers

import (
	"errors"
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
			respondWithError(w, http.StatusBadRequest, "Нет task ID")
			return
		}

		task, err := s.GetTaskByID(idStr)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				respondWithError(w, http.StatusNotFound, "Задача не найдена")
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		if task.Repeat != "" {
			nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			err = s.UpdateTaskForDone(model.TaskEntity{
				ID:      task.ID,
				Date:    nextDate,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			})
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			err = s.DeleteTask(task.ID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		respondWithJSON(w, http.StatusOK, struct{}{})
	}
}
