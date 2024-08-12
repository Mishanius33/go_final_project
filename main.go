package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mishanius33/go_final_project/common"
	"github.com/mishanius33/go_final_project/storage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := storage.NewStorage()

	// Server
	webDir := "./web"

	todoPort := os.Getenv("TODO_PORT")
	if todoPort == "" {
		todoPort = common.Port
	}

	http.HandleFunc("/api/nextdate", storage.NextDateHandler(db))
	http.HandleFunc("/api/task", storage.TaskHandler(db))
	http.HandleFunc("/api/tasks", storage.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", storage.TaskDoneHandler(db))

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":"+todoPort, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
}
