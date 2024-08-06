package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mishanius33/go_final_project/nextdate"
)

type taskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type taskResponse struct {
	ID  int64  `json:"id"`
	Err string `json:"error"`
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task taskRequest
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			sendJSONResponse(w, taskResponse{Err: "Error decoding JSON: " + err.Error()}, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			sendJSONResponse(w, taskResponse{Err: "Требуется title"}, http.StatusBadRequest)
			return
		}

		var date time.Time
		if task.Date == "" {
			date = time.Now()
		} else {
			date, err = time.Parse("20060102", task.Date)
			if err != nil {
				sendJSONResponse(w, taskResponse{Err: "Неверный формат даты"}, http.StatusBadRequest)
				return
			}
		}

		var nextDate string
		if date.Before(time.Now()) {
			if task.Repeat == "" {
				nextDate = time.Now().Format("20060102")
			} else {
				nextDateStr, err := nextdate.NextDate(time.Now(), date.Format("20060102"), task.Repeat)
				if err != nil {
					sendJSONResponse(w, taskResponse{Err: "Неверный формат повторения: " + task.Repeat}, http.StatusBadRequest)
					return
				}
				nextDate = nextDateStr
			}
		} else {
			nextDate = date.Format("20060102")
		}

		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", nextDate, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, taskResponse{ID: id}, http.StatusOK)
	}
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
