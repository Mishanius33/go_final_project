package handlers

import (
	"net/http"

	"github.com/mishanius33/go_final_project/storage"
)

func TaskHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				data, status, err := s.GetTaskByID(id)
				if err != nil {
					http.Error(w, string(data), status)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(status)

				_, err = w.Write(data)
				if err != nil {
					http.Error(w, "Ошибка ответа", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, `{"error": "Идентификатор не указан"}`, http.StatusBadRequest)
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
