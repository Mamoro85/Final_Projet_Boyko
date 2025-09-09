package api

import (
	"fmt"
	"net/http"

	"github.com/Mamoro85/Final_Projet_Boyko.git/pkg/db"
)

// TaskResp — структура ответа для отдельной задачи (с ID как строка)
type TaskResp struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TasksResp — структура ответа для списка задач
type TasksResp struct {
	Tasks []TaskResp `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаем задачи из базы данных (максимум 50)
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Конвертируем задачи в формат с ID как строка
	taskResps := make([]TaskResp, len(tasks))
	for i, task := range tasks {
		taskResps[i] = TaskResp{
			ID:      fmt.Sprintf("%d", task.ID),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}
	}

	// Возвращаем список задач
	writeJSON(w, TasksResp{
		Tasks: taskResps,
	})
}
