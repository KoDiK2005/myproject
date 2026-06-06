package handler

import (
	"myproject/internal/models"
	"myproject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// URL: /users/42 — берём последний сегмент пути
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.svc.GetUserByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input models.CreateUserInput
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.UpdateUser(id, input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
