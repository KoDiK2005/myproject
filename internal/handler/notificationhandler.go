package handler

import (
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// GetNotifications godoc
// @Summary      Лента уведомлений
// @Tags         notifications
// @Produce      json
// @Param        page  query int false "Страница" default(1)
// @Param        limit query int false "Лимит"    default(20)
// @Success      200 {object} models.PaginatedResponse
// @Security     BearerAuth
// @Router       /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	page, limit, ok := parsePagination(c, 20, 100)
	if !ok {
		return
	}
	userID := c.GetInt("user_id")
	resp, err := h.svc.GetByUserID(userID, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, resp)
}

// GetUnreadCount godoc
// @Summary      Количество непрочитанных уведомлений
// @Tags         notifications
// @Produce      json
// @Success      200 {object} map[string]int
// @Security     BearerAuth
// @Router       /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt("user_id")
	count, err := h.svc.UnreadCount(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, gin.H{"count": count})
}

// MarkRead godoc
// @Summary      Отметить уведомление прочитанным
// @Tags         notifications
// @Param        id path int true "ID уведомления"
// @Success      204
// @Security     BearerAuth
// @Router       /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.MarkRead(userID, id); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(204)
}

// MarkAllRead godoc
// @Summary      Отметить все уведомления прочитанными
// @Tags         notifications
// @Success      204
// @Security     BearerAuth
// @Router       /notifications/read-all [post]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	if err := h.svc.MarkAllRead(userID); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(204)
}
