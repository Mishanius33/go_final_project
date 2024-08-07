package handlers

import (
	"database/sql"
	"net/http"
)

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			sendJSONResponse(w, map[string]string{"error": "Нет task ID"}, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			sendJSONResponse(w, map[string]string{"error": "Не найдена задача"}, http.StatusNotFound)
			return
		}

		sendJSONResponse(w, struct{}{}, http.StatusOK)
	}
}
