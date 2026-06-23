package handler

import (
	"myproject/internal/models"
	"myproject/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshCookieName   = "refresh_token"
	refreshCookiePath   = "/auth"
	refreshCookieMaxAge = 7 * 24 * 60 * 60 // 7 дней, секунды
)

type AuthHandler struct {
	svc        *service.UserService
	refreshSvc *service.RefreshTokenService
	secret     string
}

func NewAuthHandler(svc *service.UserService, refreshSvc *service.RefreshTokenService, secret string) *AuthHandler {
	return &AuthHandler{svc: svc, refreshSvc: refreshSvc, secret: secret}
}

// setRefreshCookie — кладёт refresh-токен в httpOnly cookie, недоступную для JS (защита от XSS).
// Secure=false для локальной разработки по http; в проде за HTTPS должно быть true.
func setRefreshCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, token, refreshCookieMaxAge, refreshCookiePath, "", false, true)
}

func clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", false, true)
}

func (h *AuthHandler) makeAccessToken(userID int) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return accessToken.SignedString([]byte(h.secret))
}

// Refresh godoc
// @Summary      Обновить access token (refresh token берётся из httpOnly cookie)
// @Tags         auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTok, err := c.Cookie(refreshCookieName)
	if err != nil || refreshTok == "" {
		c.JSON(401, gin.H{"error": "refresh token required"})
		return
	}

	rt, err := h.refreshSvc.ValidateToken(refreshTok)
	if err != nil {
		clearRefreshCookie(c)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// ротация: старый токен сразу инвалидируем — повторно использовать нельзя
	h.refreshSvc.RevokeToken(refreshTok)

	accessTokenStr, err := h.makeAccessToken(rt.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	newRT, err := h.refreshSvc.GenerateToken(rt.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}
	setRefreshCookie(c, newRT.Token)

	c.JSON(200, gin.H{"access_token": accessTokenStr})
}

// Logout godoc
// @Summary      Выйти (отозвать refresh token из cookie)
// @Tags         auth
// @Success      204
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	if refreshTok, err := c.Cookie(refreshCookieName); err == nil && refreshTok != "" {
		h.refreshSvc.RevokeToken(refreshTok)
	}
	clearRefreshCookie(c)
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
	clearRefreshCookie(c)
	c.Status(204)
}

// Login godoc
// @Summary      Получить JWT токен (refresh token уходит в httpOnly cookie)
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

	accessTokenStr, err := h.makeAccessToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.refreshSvc.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}
	setRefreshCookie(c, refreshToken.Token)

	c.JSON(200, gin.H{"access_token": accessTokenStr})
}
