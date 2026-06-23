package handler

import (
	"errors"
	"myproject/internal/models"
	"myproject/internal/service"

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
	page, limit, ok := parsePagination(c, 10, 100)
	if !ok {
		return
	}

	search := c.Query("search")

	var resp *models.PaginatedResponse
	var err error
	if search != "" {
		resp, err = h.svc.SearchPosts(search, page, limit)
	} else if userID := c.GetInt("user_id"); userID != 0 {
		// авторизован — персональная лента (свои + друзья + публичные)
		resp, err = h.svc.GetFeed(userID, page, limit)
	} else {
		// гость — только публичные
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
	userID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}

	page, limit, ok := parsePagination(c, 10, 100)
	if !ok {
		return
	}

	viewerID := c.GetInt("user_id") // 0 если гость
	resp, err := h.svc.GetPostsByUserID(userID, viewerID, page, limit)
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
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}
	viewerID := c.GetInt("user_id") // 0 если гость
	post, err := h.svc.GetPostByID(id, viewerID)
	if err != nil {
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
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
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
	id, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}
	userID := c.GetInt("user_id")
	if err := h.svc.DeletePost(id, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(404, gin.H{"error": "post not found"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(403, gin.H{"error": "you can't delete someone else's post"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(204)
}
