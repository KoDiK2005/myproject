package handler

import (
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

	resp, err := h.svc.ListPosts(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, resp)
}

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
		switch err.Error() {
		case "post not found":
			c.JSON(404, gin.H{"error": "post not found"})
		case "forbidden":
			c.JSON(403, gin.H{"error": "you can't edit someone else's post"})
		default:
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, post)
}

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
