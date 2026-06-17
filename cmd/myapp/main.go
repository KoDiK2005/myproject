package main

import (
	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/logger"
	"myproject/internal/repository"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	logger.Init()

	cfg := config.Load()
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(userSvc, cfg.JWTSecret)

	postRepo := repository.NewPostRepo(db)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	gin.SetMode(gin.ReleaseMode) // убираем дефолтный дебаг-вывод Gin
	r := gin.New()
	r.Use(handler.LoggerMiddleware())
	r.Use(handler.RateLimitMiddleware())
	r.Use(gin.Recovery())
	r.POST("/auth/login", authHandler.Login)

	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.ListUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
		api.GET("/users/:id/posts", postHandler.GetPostsByUser)

		api.GET("/posts", postHandler.ListPosts)
		api.GET("/posts/:id", postHandler.GetPost)
	}
	protected := api.Group("")
	protected.Use(handler.AuthMiddleware(cfg.JWTSecret))
	{
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
	}

	// Запускаем сервер (порт можно взять из cfg.Port, например ":8080")
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	logger.Log.Info().Str("port", port).Msg("server starting")
	r.Run(":" + port)
}
