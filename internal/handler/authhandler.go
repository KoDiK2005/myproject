package handler

import (
	"myproject/internal/models"
	"myproject/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	svc    *service.UserService
	secret string
}

func NewAuthHandler(svc *service.UserService, secret string) *AuthHandler {
	return &AuthHandler{svc: svc, secret: secret}
}

// Login godoc
// @Summary      Получить JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.LoginInput true "Email и пароль"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Login(input)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenStr})

}
