package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверяем пустое правило
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	// Парсим начальную дату
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %w", err)
	}

	// Приводим now к дате без времени
	nowDate, _ := time.Parse("20060102", now.Format("20060102"))

	// Разбираем правило
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d": // формат d <число>
		if len(parts) != 2 {
			return "", errors.New("invalid format for daily repeat")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("invalid day interval, must be 1-400")
		}

		// Сдвигаем дату, пока не будет больше now
		date := startDate
		for !afterNow(date, nowDate) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format("20060102"), nil

	case "y": // ежегодное повторение
		date := startDate
		for !afterNow(date, nowDate) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format("20060102"), nil

	default:
		// неподдерживаемый формат
		return "", fmt.Errorf("unsupported repeat format: %s", repeat)
	}
}

// afterNow возвращает true, если first строго больше second (без времени)
func afterNow(first, second time.Time) bool {
	return first.After(second)
}
