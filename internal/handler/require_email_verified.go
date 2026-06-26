package handler

import (
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

// RequireEmailVerified — блокирует действие, если у пользователя не подтверждён email.
// Должен идти после AuthMiddleware (читает user_id из контекста).
func RequireEmailVerified(userSvc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")
		user, err := userSvc.GetUserByID(userID)
		if err != nil || user == nil {
			c.JSON(404, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		if !user.EmailVerified {
			c.JSON(403, gin.H{"error": "email verification required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
