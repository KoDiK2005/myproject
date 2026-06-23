package handler

import (
	"myproject/internal/models"
	"myproject/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type refreshInput struct {
	Token string `json:"refresh_token" binding:"required"`
}

type AuthHandler struct {
	svc        *service.UserService
	refreshSvc *service.RefreshTokenService
	secret     string
}

func NewAuthHandler(svc *service.UserService, refreshSvc *service.RefreshTokenService, secret string) *AuthHandler {
	return &AuthHandler{svc: svc, refreshSvc: refreshSvc, secret: secret}
}

// Refresh godoc
// @Summary      Обновить access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body refreshInput true "Refresh token"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var input refreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rt, err := h.refreshSvc.ValidateToken(input.Token)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// ротация: старый токен сразу инвалидируем — повторно использовать нельзя
	h.refreshSvc.RevokeToken(input.Token)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": rt.UserID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenStr, err := accessToken.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	// выдаём новый refresh token
	newRT, err := h.refreshSvc.GenerateToken(rt.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessTokenStr,
		"refresh_token": newRT.Token,
	})
}

// Logout godoc
// @Summary      Выйти (отозвать refresh token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body refreshInput true "Refresh token"
// @Success      204
// @Failure      400 {object} map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var input refreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.refreshSvc.RevokeToken(input.Token)
	c.Status(204)
}

// LogoutAll godoc
// @Summary      Выйти со всех устройств (отозвать все refresh-токены)
// @Tags         auth
// @Success      204
// @Security     BearerAuth
// @Router       /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID := c.GetInt("user_id")
	if err := h.refreshSvc.RevokeAllSessions(userID); err != nil {
		c.JSON(500, gin.H{"error": "failed to revoke sessions"})
		return
	}
	c.Status(204)
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

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.refreshSvc.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessTokenStr,
		"refresh_token": refreshToken.Token,
	})

}
