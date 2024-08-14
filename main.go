package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mishanius33/go_final_project/handlers"
	"github.com/mishanius33/go_final_project/storage"
)

const Port = "7540"

func main() {

	// Server
	webDir := "./web"

	todoPort := os.Getenv("TODO_PORT")
	if todoPort == "" {
		todoPort = Port
	}

	db, err := storage.NewStorage()

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler(db))
	http.HandleFunc("/api/task", handlers.TaskHandler(db))
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", handlers.TaskDoneHandler(db))

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":"+todoPort, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
}
