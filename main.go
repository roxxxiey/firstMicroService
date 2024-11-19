package main

import (
	"log"

	"controlPeople/internal/handlers"
	"controlPeople/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных
	db := repository.ConnectDB()
	defer db.Close()

	// Инициализация маршрутов
	r := gin.Default()
	handlers.RegisterRoutes(r, db)

	// Запуск сервера
	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
