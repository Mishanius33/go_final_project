package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mishanius33/go_final_project/common"
)

func (s *Storage) AddTask(t *common.TaskEntity) (int, error) {
	res, err := s.db.Exec(`INSERT INTO scheduler (date,title,comment,repeat) values (:date,:title,:comment,:repeat)`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)

	if err != nil {
		log.Println("Не добавилась задача", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Не удалось получить id", err)
		return 0, err
	}

	return int(id), nil
}

func UpdateTask(s *Storage, task common.TaskEntity) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTasks(id string) (common.TaskEntity, error) {
	var t common.TaskEntity

	row := s.db.QueryRow(`SELECT * FROM scheduler where id = :id`,
		sql.Named("id", id),
	)

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		log.Println("Не удалось получить задачу по id:", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return common.TaskEntity{}, errors.New("Задача не найдена")
		}
		return common.TaskEntity{}, err
	}

	return t, nil
}

func GetTaskByID(s *Storage, id string) ([]byte, int, error) {
	var t common.TaskEntity
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := s.db.QueryRow(query, id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return []byte(`{"error": "Задача не найдена"}`), http.StatusNotFound, nil
		}
		return nil, http.StatusInternalServerError, err
	}

	response, err := json.Marshal(t)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return response, http.StatusOK, nil
}

func DeleteTask(s *Storage, taskID string) error {
	log.Println("Удаление задачи:", taskID)

	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`,
		sql.Named("id", taskID),
	)
	if err != nil {
		log.Println("Не удалось удалить задачу:", err)
		return errors.New("Задача не найдена")
	}

	log.Println("Удаление успешно")

	return nil
}
