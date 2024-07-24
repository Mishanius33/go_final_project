package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	// DB
	// Путь к базе данных - задание со звездочкой
	TODO_DBFILE := os.Getenv("TODO_DBFILE")
	if TODO_DBFILE == "" {
		TODO_DBFILE = "scheduler.db"
	}

	appDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении директории: %v", err)
	}
	dbPath := filepath.Join(appDir, TODO_DBFILE)

	// Проверка существования дб, если нет - создаем
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {

		err = createDatabase(dbPath)
		if err != nil {
			log.Fatalf("Ошибка создания ДБ: %v", err)
		}
		fmt.Println("ДБ создана успешно.")
	} else {
		fmt.Println("ДБ уже есть.")
	}

	// Server
	webDir := "./web"

	todo_port := os.Getenv("TODO_PORT")
	if todo_port == "" {
		todo_port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":"+todo_port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}

}
