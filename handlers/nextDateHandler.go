package handlers

import (
	"database/sql"
	"go_final_project/nextdate"
	"net/http"
	"time"
)

func NextDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nowInString := r.URL.Query().Get("now")
		if nowInString == "" {
			http.Error(w, "now missing", http.StatusBadRequest)
			return
		}
		date := r.URL.Query().Get("date")
		if date == "" {
			http.Error(w, "date missing", http.StatusBadRequest)
			return
		}
		repeat := r.URL.Query().Get("repeat")
		if repeat == "" {
			http.Error(w, "repeat missing", http.StatusBadRequest)
			return
		}

		// Отбрасываем время из nowInString
		now, err := time.Parse("20060102", nowInString)
		if err != nil {
			http.Error(w, "Время не может быть преобразовано в корректную дату", http.StatusBadRequest)
			return
		}
		now = now.Truncate(24 * time.Hour)

		nextDate, err := nextdate.NextDate(now, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(nextDate))
	}
}
