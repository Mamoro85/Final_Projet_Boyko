package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная переменная для хранения соединения с БД
var DB *sql.DB

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

// Init инициализирует базу данных и создаёт таблицу при отсутствии файла
func Init(dbFile string) error {
	// Проверяем, существует ли база данных
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем соединение с SQLite
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	// Создаём схему, если файл базы данных только что создан
	if install {
		fmt.Println("Создаём базу данных и таблицу scheduler...")
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("не удалось создать схему базы данных: %w", err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
