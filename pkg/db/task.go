package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Task — описание задачи
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Tasks возвращает список задач, отсортированных по дате в порядке возрастания
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date, id LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var comment, repeat sql.NullString

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &comment, &repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		// Обрабатываем nullable поля
		if comment.Valid {
			task.Comment = comment.String
		}
		if repeat.Valid {
			task.Repeat = repeat.String
		}

		tasks = append(tasks, &task)
	}

	// Проверяем ошибки итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	// Если задач нет, возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по указанному ID
func GetTask(id string) (*Task, error) {
	var task Task
	var comment, repeat sql.NullString

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &comment, &repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}

	// Обрабатываем nullable поля
	if comment.Valid {
		task.Comment = comment.String
	}
	if repeat.Valid {
		task.Repeat = repeat.String
	}

	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества обновленных записей: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по указанному ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	// Проверяем, была ли удалена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества удаленных записей: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// UpdateDate обновляет дату задачи по указанному ID
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества обновленных записей: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
