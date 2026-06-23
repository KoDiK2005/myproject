package handler

import (
	"errors"
	"fmt"
	"io"
	"myproject/internal/logger"
	"myproject/internal/models"
	"myproject/internal/service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc            *service.UserService
	emailVerifySvc *service.EmailVerificationService
	uploadPath     string // папка для аватаров, например "./uploads/avatars"
}

func NewUserHandler(svc *service.UserService, emailVerifySvc *service.EmailVerificationService) *UserHandler {
	return &UserHandler{svc: svc, emailVerifySvc: emailVerifySvc, uploadPath: "./uploads/avatars"}
}

// ListUsers godoc
// @Summary      Список пользователей
// @Tags         users
// @Produce      json
// @Param        page  query int false "Страница" default(1)
// @Param        limit query int false "Лимит"    default(10)
// @Success      200 {object} models.PaginatedResponse
// @Failure      500 {object} map[string]string
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, limit, ok := parsePagination(c, 10, 100)
	if !ok {
		return
	}

	search := c.Query("search")
	var resp *models.PaginatedResponse
	var err error
	if search != "" {
		resp, err = h.svc.SearchUsers(search, page, limit)
	} else {
		resp, err = h.svc.ListUsers(page, limit)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, resp)
}
// GetUser godoc
// @Summary      Получить пользователя по ID
// @Tags         users
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      200 {object} models.User
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// URL: /users/42 — берём последний сегмент пути
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}

	user, err := h.svc.GetUserByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, user)
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body models.CreateUserInput true "Данные пользователя"
// @Success      201 {object} models.User
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(input)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(409, gin.H{"error": "email already taken"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	// письмо подтверждения — best-effort, регистрацию не роняем при сбое отправки
	if err := h.emailVerifySvc.SendVerification(user.ID, user.Email); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", user.ID).Msg("failed to send verification email")
	}

	c.JSON(201, user)
}

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Tags         users
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      204
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}

	if !requireOwnResource(c, id, "you can't delete someone else's account") {
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(404, gin.H{"error": "user not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(204)
}

// UpdateUser godoc
// @Summary      Обновить пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path int                    true "ID пользователя"
// @Param        input body models.UpdateUserInput true "Новые данные"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input models.UpdateUserInput
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}

	if !requireOwnResource(c, id, "you can't update someone else's account") {
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.UpdateUser(id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(404, gin.H{"error": "user not found"})
		case errors.Is(err, service.ErrEmailTaken):
			c.JSON(409, gin.H{"error": "email already taken"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, user)
}

// UploadAvatar godoc
// @Summary      Загрузить аватар пользователя
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path int  true "ID пользователя"
// @Param        avatar formData file true "Файл аватара (jpg/png, макс 2MB)"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Security     BearerAuth
// @Router       /users/{id}/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}

	if !requireOwnResource(c, id, "not your account") {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"error": "no file provided"})
		return
	}

	// 2MB максимум
	if file.Size > 2<<20 {
		c.JSON(400, gin.H{"error": "file too large (max 2MB)"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(400, gin.H{"error": "only jpg/png allowed"})
		return
	}

	// проверяем реальный MIME по magic bytes — расширение можно подделать
	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := io.ReadFull(f, buf); err != nil && err != io.ErrUnexpectedEOF {
		c.JSON(500, gin.H{"error": "failed to read file"})
		return
	}
	mime := http.DetectContentType(buf)
	if mime != "image/jpeg" && mime != "image/png" {
		c.JSON(400, gin.H{"error": "file content must be image/jpeg or image/png"})
		return
	}

	// имя файла — user_id + расширение, перезаписываем при повторной загрузке
	filename := fmt.Sprintf("%d%s", id, ext)
	dst := filepath.Join(h.uploadPath, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": "failed to save file"})
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	if err := h.svc.UpdateAvatar(id, avatarURL); err != nil {
		c.JSON(500, gin.H{"error": "failed to update avatar"})
		return
	}

	c.JSON(200, gin.H{"avatar": avatarURL})
}
