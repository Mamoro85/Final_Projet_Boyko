package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mamoro85/Final_Projet_Boyko.git/pkg/api"
	"github.com/Mamoro85/Final_Projet_Boyko.git/pkg/db"
)

func main() {
	// 1. Инициализация базы данных
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer cleanupDB()

	// 2. Инициализация API
	api.Init()

	// 3. Настройка веб-сервера
	webDir := "./web"
	port := getPort()

	fmt.Printf("Сервер запущен: http://localhost:%s\n", port)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// 4. Запуск сервера
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// cleanupDB закрывает соединение с базой данных
func cleanupDB() {
	if err := db.Close(); err != nil {
		log.Printf("Ошибка закрытия базы данных: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7541"
	}
	return port
}
