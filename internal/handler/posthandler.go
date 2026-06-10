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
	posts, err := h.svc.ListPosts()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	post, err := h.svc.GetPostByID(id)
	if err != nil {
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

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeletePost(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
