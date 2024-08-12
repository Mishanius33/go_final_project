package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mishanius33/go_final_project/common"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	storage := &Storage{}

	err := storage.initDB()

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initDB() error {
	var err error

	err = CreateDatabase(common.DbFile)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", common.DbFile)

	if err != nil {
		return fmt.Errorf("can't open db %s: %w", common.DbFile, err)
	}

	err = CreateTableAndIdx(db)

	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func CreateDatabase(dbPath string) error {
	todoDBFile := os.Getenv("TODO_DBFILE")
	if todoDBFile == "" {
		todoDBFile = common.DbFile
	}

	appPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Не удалось получить путь к исполняемому файлу: %w", err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), todoDBFile)

	_, err = os.Stat(dbFile)
	if err == nil {
		// ДБ уже существует
		fmt.Println("ДБ уже есть.")
		return nil
	} else if os.IsNotExist(err) {
		// ДБ не существует, нужно создать
		_, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("Ошибка создания ДБ: %w", err)
		}
		fmt.Println("ДБ создана успешно.")
		return nil
	} else {
		// Произошла другая ошибка
		return fmt.Errorf("Не удалось проверить наличие базы данных: %w", err)
	}
}

func CreateTableAndIdx(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  date CHAR(8) NOT NULL,
		  title VARCHAR(128) NOT NULL DEFAULT '',
		  comment TEXT DEFAULT '',
		  repeat VARCHAR(128) NOT NULL
		  );
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	  `)
	if err != nil {
		return fmt.Errorf("Не удалось создать новую ДБ %w", err)
	}

	return nil
}
