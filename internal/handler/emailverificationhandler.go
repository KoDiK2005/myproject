package handler

import (
	"errors"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type EmailVerificationHandler struct {
	svc     *service.EmailVerificationService
	userSvc *service.UserService
}

func NewEmailVerificationHandler(svc *service.EmailVerificationService, userSvc *service.UserService) *EmailVerificationHandler {
	return &EmailVerificationHandler{svc: svc, userSvc: userSvc}
}

// VerifyEmail godoc
// @Summary      Подтвердить email по токену из письма
// @Tags         auth
// @Produce      json
// @Param        token query string true "Токен из письма"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/verify-email [get]
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "token required"})
		return
	}
	if err := h.svc.VerifyToken(token); err != nil {
		if errors.Is(err, service.ErrTokenInvalid) {
			c.JSON(400, gin.H{"error": "invalid or expired token"})
			return
		}
		c.JSON(500, gin.H{"error": "verification failed"})
		return
	}
	c.JSON(200, gin.H{"status": "email verified"})
}

// ResendVerification godoc
// @Summary      Отправить письмо подтверждения повторно
// @Tags         auth
// @Success      204
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /auth/resend-verification [post]
func (h *EmailVerificationHandler) ResendVerification(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := h.userSvc.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if user.EmailVerified {
		c.JSON(400, gin.H{"error": "email already verified"})
		return
	}
	if err := h.svc.SendVerification(user.ID, user.Email); err != nil {
		c.JSON(500, gin.H{"error": "failed to send verification email"})
		return
	}
	c.Status(204)
}
