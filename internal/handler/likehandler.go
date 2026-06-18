package handler

import (
	"myproject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	svc *service.LikeService
}

func NewLikeHandler(svc *service.LikeService) *LikeHandler {
	return &LikeHandler{svc: svc}
}

// LikePost godoc
// @Summary      Лайкнуть пост
// @Tags         likes
// @Produce      json
// @Param        id path int true "ID поста"
// @Success      204
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts/{id}/like [post]
func (h *LikeHandler) LikePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid post id"})
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.Like(userID, postID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// UnlikePost godoc
// @Summary      Убрать лайк
// @Tags         likes
// @Produce      json
// @Param        id path int true "ID поста"
// @Success      204
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts/{id}/like [delete]
func (h *LikeHandler) UnlikePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid post id"})
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.Unlike(userID, postID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// GetLikeCount godoc
// @Summary      Количество лайков у поста
// @Tags         likes
// @Produce      json
// @Param        id path int true "ID поста"
// @Success      200 {object} map[string]int
// @Failure      400 {object} map[string]string
// @Router       /posts/{id}/likes [get]
func (h *LikeHandler) GetLikeCount(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid post id"})
		return
	}
	count, err := h.svc.Count(postID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"count": count})
}
