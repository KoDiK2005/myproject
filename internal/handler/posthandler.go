package handler

import (
	"errors"
	"myproject/internal/models"
	"myproject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// ListPosts godoc
// @Summary      Список постов (с опциональным поиском)
// @Tags         posts
// @Produce      json
// @Param        page   query int    false "Страница" default(1)
// @Param        limit  query int    false "Лимит"    default(10)
// @Param        search query string false "Поиск по title/body"
// @Success      200 {object} models.PaginatedResponse
// @Failure      500 {object} map[string]string
// @Router       /posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	search := c.Query("search")

	var resp *models.PaginatedResponse
	if search != "" {
		resp, err = h.svc.SearchPosts(search, page, limit)
	} else {
		resp, err = h.svc.ListPosts(page, limit)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, resp)
}

// GetPostsByUser godoc
// @Summary      Посты конкретного пользователя
// @Tags         posts
// @Produce      json
// @Param        id    path int false "ID пользователя"
// @Param        page  query int false "Страница" default(1)
// @Param        limit query int false "Лимит"    default(10)
// @Success      200 {object} models.PaginatedResponse
// @Failure      500 {object} map[string]string
// @Router       /users/{id}/posts [get]
func (h *PostHandler) GetPostsByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	resp, err := h.svc.GetPostsByUserID(userID, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, resp)
}

// GetPost godoc
// @Summary      Получить пост по ID
// @Tags         posts
// @Produce      json
// @Param        id path int true "ID поста"
// @Success      200 {object} models.Post
// @Failure      404 {object} map[string]string
// @Router       /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	post, err := h.svc.GetPostByID(id)
	if err != nil || post == nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}
	c.JSON(200, post)
}

// CreatePost godoc
// @Summary      Создать пост
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        input body models.CreatePostInput true "Данные поста"
// @Success      201 {object} models.Post
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var input models.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	post, err := h.svc.CreatePost(input, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, post)
}

// UpdatePost godoc
// @Summary      Обновить пост
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path int                    true "ID поста"
// @Param        input body models.UpdatePostInput true "Новые данные"
// @Success      200 {object} models.Post
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var input models.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("user_id")
	post, err := h.svc.UpdatePost(id, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(404, gin.H{"error": "post not found"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(403, gin.H{"error": "you can't edit someone else's post"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, post)
}

// DeletePost godoc
// @Summary      Удалить пост
// @Tags         posts
// @Produce      json
// @Param        id path int true "ID поста"
// @Success      204
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.DeletePost(id, userID); err != nil {
		switch err.Error() {
		case "post not found":
			c.JSON(404, gin.H{"error": "post not found"})
		case "forbidden":
			c.JSON(403, gin.H{"error": "you can't delete someone else's post"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(204)
}
