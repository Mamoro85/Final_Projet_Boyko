package api

import (
	"net/http"

	"github.com/Mamoro85/Final_Projet_Boyko.git/pkg/db"
)

// TasksResp — структура ответа для списка задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
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

	// Возвращаем список задач
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
