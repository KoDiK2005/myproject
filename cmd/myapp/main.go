// @title           myproject API
// @version         1.0
// @description     Учебный REST API на Go + Gin + PostgreSQL
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "myproject/docs"
	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/logger"
	"myproject/internal/repository"
	"myproject/internal/service"
	"myproject/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("invalid configuration")
	}
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// connection pool: без этого при нагрузке база захлебнётся
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	userRepo := repository.NewUserRepo(db)
	postRepo := repository.NewPostRepo(db)
	refreshTokenRepo := repository.NewRefreshTokenRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	likeRepo := repository.NewLikeRepo(db)
	friendshipRepo := repository.NewFriendshipRepo(db)
	messageRepo := repository.NewMessageRepo(db)

	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)
	refreshTokenSvc := service.NewRefreshTokenService(refreshTokenRepo)
	commentSvc := service.NewCommentService(commentRepo, postSvc)
	likeSvc := service.NewLikeService(likeRepo)
	friendshipSvc := service.NewFriendshipService(friendshipRepo)
	messageSvc := service.NewMessageService(messageRepo)

	wsHub := ws.NewHub()
	go wsHub.Run() // запускаем горутину-диспетчер

	userHandler := handler.NewUserHandler(userSvc)
	postHandler := handler.NewPostHandler(postSvc)
	authHandler := handler.NewAuthHandler(userSvc, refreshTokenSvc, cfg.JWTSecret)
	commentHandler := handler.NewCommentHandler(commentSvc)
	likeHandler := handler.NewLikeHandler(likeSvc)
	friendshipHandler := handler.NewFriendshipHandler(friendshipSvc)
	messageHandler := handler.NewMessageHandler(messageSvc, wsHub, cfg.AllowedOrigins)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(handler.LoggerMiddleware())
	r.Use(handler.PrometheusMiddleware())
	r.Use(handler.RateLimitMiddleware())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			c.JSON(503, gin.H{"status": "unavailable", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.Refresh)
	r.POST("/auth/logout", authHandler.Logout)

	// WebSocket — авторизация через query param ?token=...
	r.GET("/ws", handler.AuthMiddleware(cfg.JWTSecret), messageHandler.WSConnect)

	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.ListUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
		api.GET("/users/:id/posts", handler.OptionalAuthMiddleware(cfg.JWTSecret), postHandler.GetPostsByUser)

		// лента с optional auth — авторизованные получают персональный фид
		api.GET("/posts", handler.OptionalAuthMiddleware(cfg.JWTSecret), postHandler.ListPosts)
		api.GET("/posts/:id", handler.OptionalAuthMiddleware(cfg.JWTSecret), postHandler.GetPost)
		api.GET("/posts/:id/comments", handler.OptionalAuthMiddleware(cfg.JWTSecret), commentHandler.ListComments)
		api.GET("/posts/:id/likes", handler.OptionalAuthMiddleware(cfg.JWTSecret), likeHandler.GetLikeCount)
	}

	protected := api.Group("")
	protected.Use(handler.AuthMiddleware(cfg.JWTSecret))
	{
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.POST("/users/:id/avatar", userHandler.UploadAvatar)

		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)

		protected.POST("/posts/:id/comments", commentHandler.CreateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)

		protected.POST("/posts/:id/like", likeHandler.LikePost)
		protected.DELETE("/posts/:id/like", likeHandler.UnlikePost)

		// мессенджер
		protected.GET("/messages", messageHandler.GetConversations)
		protected.GET("/messages/:id", messageHandler.GetHistory)
		protected.POST("/messages/:id", messageHandler.SendMessage)

		// дружба
		protected.POST("/friends/request/:id", friendshipHandler.SendRequest)
		protected.POST("/friends/accept/:id", friendshipHandler.Accept)
		protected.POST("/friends/reject/:id", friendshipHandler.Reject)
		protected.DELETE("/friends/:id", friendshipHandler.Remove)
		protected.GET("/friends", friendshipHandler.GetFriends)
		protected.GET("/friends/requests/incoming", friendshipHandler.GetIncomingRequests)
		protected.GET("/friends/requests/outgoing", friendshipHandler.GetOutgoingRequests)
		protected.GET("/friends/status/:id", friendshipHandler.GetStatus)
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.Log.Info().Str("port", port).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("forced shutdown")
	}

	logger.Log.Info().Msg("server stopped")
}
