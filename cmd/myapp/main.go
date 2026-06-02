package main

import (
	"log"
	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/repository"
	"myproject/internal/service"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Регистрируем маршрут для создания пользователя
	http.HandleFunc("/users", userHandler.CreateUser) // POST /users

	// Запускаем сервер (порт можно взять из cfg.Port, например ":8080")
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
