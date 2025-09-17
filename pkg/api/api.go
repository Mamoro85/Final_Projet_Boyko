package api

import (
	"fmt"
	"net/http"
	"time"
)

// Init регистрирует API-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/", taskByIDHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
}

// NextDateHandler — HTTP-обработчик для /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметры
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dstart == "" || repeat == "" {
		writeJSON(w, map[string]string{"error": "missing parameters: date and repeat are required"}, http.StatusBadRequest)
		return
	}

	// если now не передан, используем текущую дату
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSON(w, map[string]string{"error": fmt.Sprintf("invalid now date: %v", err)}, http.StatusBadRequest)
			return
		}
	}

	// вычисляем следующую дату
	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Возвращаем просто строку, как ожидает тест
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
