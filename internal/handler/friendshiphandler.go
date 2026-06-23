package handler

import (
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type FriendshipHandler struct {
	svc *service.FriendshipService
}

func NewFriendshipHandler(svc *service.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{svc: svc}
}

// SendRequest godoc
// @Summary      Отправить заявку в друзья
// @Tags         friends
// @Param        id path int true "ID пользователя"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /friends/request/{id} [post]
func (h *FriendshipHandler) SendRequest(c *gin.Context) {
	targetID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.SendRequest(myID, targetID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "request sent"})
}

// Accept godoc
// @Summary      Принять заявку в друзья
// @Tags         friends
// @Param        id path int true "ID пользователя который отправил заявку"
// @Success      200 {object} map[string]string
// @Security     BearerAuth
// @Router       /friends/accept/{id} [post]
func (h *FriendshipHandler) Accept(c *gin.Context) {
	requesterID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.Accept(myID, requesterID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "accepted"})
}

// Reject godoc
// @Summary      Отклонить заявку в друзья
// @Tags         friends
// @Param        id path int true "ID пользователя который отправил заявку"
// @Success      200 {object} map[string]string
// @Security     BearerAuth
// @Router       /friends/reject/{id} [post]
func (h *FriendshipHandler) Reject(c *gin.Context) {
	requesterID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.Reject(myID, requesterID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "rejected"})
}

// Remove godoc
// @Summary      Удалить из друзей
// @Tags         friends
// @Param        id path int true "ID друга"
// @Success      204
// @Security     BearerAuth
// @Router       /friends/{id} [delete]
func (h *FriendshipHandler) Remove(c *gin.Context) {
	friendID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.Remove(myID, friendID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// GetFriends godoc
// @Summary      Список друзей
// @Tags         friends
// @Produce      json
// @Success      200 {array} models.User
// @Security     BearerAuth
// @Router       /friends [get]
func (h *FriendshipHandler) GetFriends(c *gin.Context) {
	myID := c.GetInt("user_id")
	friends, err := h.svc.GetFriends(myID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, friends)
}

// GetIncomingRequests godoc
// @Summary      Входящие заявки в друзья
// @Tags         friends
// @Produce      json
// @Success      200 {array} models.User
// @Security     BearerAuth
// @Router       /friends/requests/incoming [get]
func (h *FriendshipHandler) GetIncomingRequests(c *gin.Context) {
	myID := c.GetInt("user_id")
	users, err := h.svc.GetIncomingRequests(myID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// GetOutgoingRequests godoc
// @Summary      Исходящие заявки в друзья
// @Tags         friends
// @Produce      json
// @Success      200 {array} models.User
// @Security     BearerAuth
// @Router       /friends/requests/outgoing [get]
func (h *FriendshipHandler) GetOutgoingRequests(c *gin.Context) {
	myID := c.GetInt("user_id")
	users, err := h.svc.GetOutgoingRequests(myID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// GetStatus godoc
// @Summary      Статус дружбы с конкретным пользователем
// @Tags         friends
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      200 {object} models.FriendshipStatus
// @Security     BearerAuth
// @Router       /friends/status/{id} [get]
func (h *FriendshipHandler) GetStatus(c *gin.Context) {
	otherID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	status, err := h.svc.GetStatus(myID, otherID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, status)
}
