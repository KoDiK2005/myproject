package main

import (
	"log"
	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/repository"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"

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

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.ListUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Запускаем сервер (порт можно взять из cfg.Port, например ":8080")
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
