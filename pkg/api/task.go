package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Mamoro85/Final_Projet_Boyko.git/pkg/db"
)

// checkAndNormalizeDate validates date and repeat; ensures date is not in the past
func checkAndNormalizeDate(date, title, comment, repeat string) (string, string, string, string, error) {
	if strings.TrimSpace(title) == "" {
		return date, title, comment, repeat, fmt.Errorf("Не указан заголовок задачи")
	}
	now := time.Now()
	if strings.TrimSpace(date) == "" {
		date = now.Format(DateFormat)
	}
	parsed, err := time.Parse(DateFormat, date)
	if err != nil {
		return date, title, comment, repeat, fmt.Errorf("invalid date format")
	}
	if strings.TrimSpace(repeat) != "" {
		if _, err := NextDate(now, date, repeat); err != nil {
			return date, title, comment, repeat, fmt.Errorf("invalid repeat rule: %w", err)
		}
	}
	if !afterNow(parsed, now) { // date is today or past → fix
		if strings.TrimSpace(repeat) == "" {
			date = now.Format(DateFormat)
		} else {
			next, _ := NextDate(now, date, repeat)
			date = next
		}
	}
	return date, title, comment, repeat, nil
}

// taskHandler handles POST(add), GET(fetch by id), PUT(update), DELETE(delete)
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var in struct {
			Date    string `json:"date"`
			Title   string `json:"title"`
			Comment string `json:"comment"`
			Repeat  string `json:"repeat"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil && err.Error() != "EOF" {
			writeJSON(w, map[string]string{"error": "invalid JSON"}, http.StatusBadRequest)
			return
		}
		date, title, comment, repeat, err := checkAndNormalizeDate(in.Date, in.Title, in.Comment, in.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		res, err := db.DB.Exec(`INSERT INTO scheduler(date,title,comment,repeat) VALUES(?,?,?,?)`, date, title, comment, repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusOK)

	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
			return
		}

		task, err := db.GetTask(idStr)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			return
		}

		writeJSON(w, map[string]string{
			"id":      fmt.Sprintf("%d", task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}, http.StatusOK)

	case http.MethodPut:
		var in map[string]any
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, map[string]string{"error": "invalid JSON"}, http.StatusBadRequest)
			return
		}

		idStr := fmt.Sprint(in["id"])
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			writeJSON(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
			return
		}

		date := fmt.Sprint(in["date"])
		title := fmt.Sprint(in["title"])
		comment := fmt.Sprint(in["comment"])
		repeat := fmt.Sprint(in["repeat"])

		date, title, comment, repeat, err = checkAndNormalizeDate(date, title, comment, repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		task := &db.Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}

		err = db.UpdateTask(task)
		if err != nil {
			// Проверяем, является ли ошибка "задача не найдена"
			if err.Error() == "Задача не найдена" {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, map[string]any{}, http.StatusOK)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
			return
		}

		err := db.DeleteTask(idStr)
		if err != nil {
			// Проверяем, является ли ошибка "задача не найдена"
			if err.Error() == "Задача не найдена" {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, map[string]any{}, http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Получаем задачу по ID
	task, err := db.GetTask(idStr)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	// Если задача не повторяющаяся, удаляем её
	if strings.TrimSpace(task.Repeat) == "" {
		err = db.DeleteTask(idStr)
		if err != nil {
			// Проверяем, является ли ошибка "задача не найдена"
			if err.Error() == "Задача не найдена" {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			}
			return
		}
		writeJSON(w, map[string]any{}, http.StatusOK)
		return
	}

	// Для повторяющейся задачи вычисляем следующую дату
	start, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка парсинга даты"}, http.StatusBadRequest)
		return
	}

	next, err := NextDate(start, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновляем дату задачи
	err = db.UpdateDate(next, idStr)
	if err != nil {
		// Проверяем, является ли ошибка "задача не найдена"
		if err.Error() == "Задача не найдена" {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]any{}, http.StatusOK)
}

// taskByIDHandler handles GET, PUT, DELETE for /api/task/{id}
func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/task/") {
		writeJSON(w, map[string]string{"error": "Invalid path"}, http.StatusBadRequest)
		return
	}

	idStr := strings.TrimPrefix(path, "/api/task/")
	if idStr == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Validate ID format
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, map[string]string{"error": "Неверный идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := db.GetTask(idStr)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			return
		}

		writeJSON(w, map[string]string{
			"id":      fmt.Sprintf("%d", task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}, http.StatusOK)

	case http.MethodPut:
		var in struct {
			Date    string `json:"date"`
			Title   string `json:"title"`
			Comment string `json:"comment"`
			Repeat  string `json:"repeat"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, map[string]string{"error": "invalid JSON"}, http.StatusBadRequest)
			return
		}

		date, title, comment, repeat, err := checkAndNormalizeDate(in.Date, in.Title, in.Comment, in.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		task := &db.Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}

		err = db.UpdateTask(task)
		if err != nil {
			// Проверяем, является ли ошибка "задача не найдена"
			if err.Error() == "Задача не найдена" {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, map[string]any{}, http.StatusOK)

	case http.MethodDelete:
		err := db.DeleteTask(idStr)
		if err != nil {
			// Проверяем, является ли ошибка "задача не найдена"
			if err.Error() == "Задача не найдена" {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, map[string]any{}, http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
