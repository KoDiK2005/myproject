package handler

import (
	"errors"
	"myproject/internal/models"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// ListComments godoc
// @Summary      Список комментариев к посту
// @Tags         comments
// @Produce      json
// @Param        id    path int true  "ID поста"
// @Param        page  query int false "Страница" default(1)
// @Param        limit query int false "Лимит"    default(20)
// @Success      200 {object} models.PaginatedResponse
// @Failure      400 {object} map[string]string
// @Router       /posts/{id}/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	postID, ok := parseIDParam(c, "id", "invalid post id")
	if !ok {
		return
	}
	page, limit, ok := parsePagination(c, 20, 100)
	if !ok {
		return
	}
	viewerID := c.GetInt("user_id") // 0 если гость
	resp, err := h.svc.ListComments(postID, viewerID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(404, gin.H{"error": "post not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

// CreateComment godoc
// @Summary      Добавить комментарий
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id    path int                       true "ID поста"
// @Param        input body models.CreateCommentInput true "Текст комментария"
// @Success      201 {object} models.Comment
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID, ok := parseIDParam(c, "id", "invalid post id")
	if !ok {
		return
	}
	var input models.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("user_id")
	comment, err := h.svc.AddComment(postID, userID, input)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(404, gin.H{"error": "post not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, comment)
}

// DeleteComment godoc
// @Summary      Удалить комментарий
// @Tags         comments
// @Produce      json
// @Param        id path int true "ID комментария"
// @Success      204
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.DeleteComment(id, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(404, gin.H{"error": "comment not found"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(403, gin.H{"error": "not your comment"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(204)
}
