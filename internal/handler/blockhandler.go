package handler

import (
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type BlockHandler struct {
	svc *service.BlockService
}

func NewBlockHandler(svc *service.BlockService) *BlockHandler {
	return &BlockHandler{svc: svc}
}

// Block godoc
// @Summary      Заблокировать пользователя
// @Tags         blocks
// @Param        id path int true "ID пользователя"
// @Success      204
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /blocks/{id} [post]
func (h *BlockHandler) Block(c *gin.Context) {
	targetID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.Block(myID, targetID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// Unblock godoc
// @Summary      Разблокировать пользователя
// @Tags         blocks
// @Param        id path int true "ID пользователя"
// @Success      204
// @Security     BearerAuth
// @Router       /blocks/{id} [delete]
func (h *BlockHandler) Unblock(c *gin.Context) {
	targetID, ok := parseIDParam(c, "id", "invalid user id")
	if !ok {
		return
	}
	myID := c.GetInt("user_id")
	if err := h.svc.Unblock(myID, targetID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// GetBlockedUsers godoc
// @Summary      Список заблокированных пользователей
// @Tags         blocks
// @Produce      json
// @Success      200 {array} models.User
// @Security     BearerAuth
// @Router       /blocks [get]
func (h *BlockHandler) GetBlockedUsers(c *gin.Context) {
	myID := c.GetInt("user_id")
	users, err := h.svc.GetBlockedUsers(myID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}
