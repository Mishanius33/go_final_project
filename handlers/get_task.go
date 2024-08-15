package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mishanius33/go_final_project/model"
	"github.com/mishanius33/go_final_project/storage"
)

func GetTasksHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := s.GetList()
		if err != nil {
			log.Printf("Не удалось получить задачи в запросе: %v", err)
			http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
			return
		}

		if len(tasks) == 0 {
			tasks = []model.TaskEntity{}
		}

		w.Header().Set("Content-Type", "application/json")
		resp := map[string][]model.TaskEntity{
			"tasks": tasks,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Ошибка сериализации JSON: %v", err)
			http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
			return
		}

		log.Printf("Успешный вывод задач. Задач найдено: %d", len(tasks))
	}
}
