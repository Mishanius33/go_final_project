package storage

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func DeleteTaskHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id не число:", err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err = DeleteTask(s, id)
		if err != nil {
			log.Println("Не удалось удалить задачу")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err = json.NewEncoder(w).Encode(map[string]string{})
		if err != nil {
			log.Println("err encode:", err)
			http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
		}
	}
}
