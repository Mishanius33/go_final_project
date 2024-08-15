package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mishanius33/go_final_project/storage"
)

func DeleteTaskHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id не число:", err)
			respondWithError(w, http.StatusBadRequest, "id не число")
			return
		}

		err = s.DeleteTask(id)
		if err != nil {
			log.Println("Не удалось удалить задачу:", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{})
	}
}
