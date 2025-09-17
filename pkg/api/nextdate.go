package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// afterNow проверяет, больше ли дата, чем now (игнорируем время суток)
func afterNow(date, now time.Time) bool {
	d := date.Truncate(24 * time.Hour)
	n := now.Truncate(24 * time.Hour)
	return d.After(n) || d.Equal(n)
}

// NextDate вычисляет следующую дату задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	// парсим стартовую дату
	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart: %w", err)
	}

	// разбираем правило
	parts := strings.Split(strings.TrimSpace(repeat), " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat format for days")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("invalid day interval")
		}
		date := start
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat format for yearly")
		}
		date := start
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format: %s", repeat)
	}
}
