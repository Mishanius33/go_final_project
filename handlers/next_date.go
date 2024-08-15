package handlers

import (
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/nextdate"
	"github.com/mishanius33/go_final_project/storage"
)

const (
	HoursPerDay = 24
	DateFormat  = "20060102"
)

func NextDateHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nowInString := r.URL.Query().Get("now")
		if nowInString == "" {
			http.Error(w, "Отсутствует now", http.StatusBadRequest)
			return
		}
		date := r.URL.Query().Get("date")
		if date == "" {
			http.Error(w, "Отсутствует date", http.StatusBadRequest)
			return
		}
		repeat := r.URL.Query().Get("repeat")
		if repeat == "" {
			http.Error(w, "Отсутствует repeat", http.StatusBadRequest)
			return
		}

		// Отбрасываем время
		now, err := time.Parse(DateFormat, nowInString)
		if err != nil {
			http.Error(w, "Время не может быть преобразовано в корректную дату", http.StatusBadRequest)
			return
		}
		now = now.Truncate(HoursPerDay * time.Hour)

		nextDate, err := nextdate.NextDate(now, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(nextDate))
		if err != nil {
			http.Error(w, "Ошибка ответа", http.StatusInternalServerError)
			return
		}
	}
}
