package handlers

import (
	"errors"
	"net/http"

	"github.com/mishanius33/go_final_project/storage"
)

func TaskHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				task, err := s.GetTaskByID(id)
				if err != nil {
					if errors.Is(err, ErrTaskNotFound) {
						respondWithError(w, http.StatusNotFound, "Задача не найдена")
					} else {
						respondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
					}
					return
				}

				respondWithJSON(w, http.StatusOK, task)
			} else {
				respondWithError(w, http.StatusBadRequest, "Идентификатор не указан")
			}
		case http.MethodPost:
			AddTaskHandler(s)(w, r)
		case http.MethodPut:
			EditTaskHandler(s)(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(s)(w, r)
		}
	}
}
