package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mishanius33/go_final_project/common"
)

func GetTasksHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := s.GetList()
		if err != nil {
			log.Printf("Не удалось получить задачи в запросе: %v", err)
			http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
			return
		}

		if len(tasks) == 0 {
			tasks = []common.TaskEntity{}
		}

		w.Header().Set("Content-Type", "application/json")
		resp := map[string][]common.TaskEntity{
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

func (s *Storage) GetList() ([]common.TaskEntity, error) {
	var limit = common.Limit
	rows, err := s.db.Query(`SELECT * FROM scheduler ORDER BY date limit :limit`, sql.Named("limit", limit))
	if err != nil {
		log.Println("Не удалось получить задачи в запросе", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []common.TaskEntity
	for rows.Next() {
		t := common.TaskEntity{}
		err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("Не удалось получить задачи в запросе", err)
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		log.Println("Не удалось получить задачи в запросе", rows.Err())
		return nil, rows.Err()
	}

	return tasks, nil
}
