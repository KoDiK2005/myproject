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
	authHandler := handler.NewAuthHandler(userSvc, cfg.JWTSecret)

	postRepo := repository.NewPostRepo(db)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	r := gin.Default()
	r.POST("/auth/login", authHandler.Login)

	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.ListUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)

		api.GET("/posts", postHandler.ListPosts)
		api.GET("/posts/:id", postHandler.GetPost)
	}
	protected := api.Group("")
	protected.Use(handler.AuthMiddleware(cfg.JWTSecret))
	{
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.POST("/posts", postHandler.CreatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
	}

	// Запускаем сервер (порт можно взять из cfg.Port, например ":8080")
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
