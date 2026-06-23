package handler

import (
	"errors"
	"myproject/internal/logger"
	"myproject/internal/models"
	"myproject/internal/service"
	"myproject/internal/ws"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	svc      *service.MessageService
	hub      *ws.Hub
	upgrader websocket.Upgrader
}

func NewMessageHandler(svc *service.MessageService, hub *ws.Hub, allowedOrigins []string) *MessageHandler {
	return &MessageHandler{
		svc: svc,
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// разрешаем с нашего фронт-домена
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				for _, allowed := range allowedOrigins {
					if origin == allowed {
						return true
					}
				}
				return false
			},
		},
	}
}

// WSConnect godoc
// @Summary      WebSocket подключение для получения сообщений
// @Tags         messages
// @Router       /ws [get]
// @Security     BearerAuth
func (h *MessageHandler) WSConnect(c *gin.Context) {
	userID := c.GetInt("user_id")
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return // upgrader уже написал 400
	}
	ws.NewClient(h.hub, userID, conn)
}

// SendMessage godoc
// @Summary      Отправить сообщение другу
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID получателя"
// @Param        body body object true "Текст"
// @Success      201 {object} models.Message
// @Failure      400,403,500 {object} map[string]string
// @Security     BearerAuth
// @Router       /messages/{id} [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	senderID := c.GetInt("user_id")
	receiverID, ok := parseIDParam(c, "id", "invalid receiver id")
	if !ok {
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=1,max=5000"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.svc.Send(senderID, receiverID, input.Content)
	if err != nil {
		if errors.Is(err, service.ErrNotFriends) {
			c.JSON(403, gin.H{"error": "you can only message friends"})
			return
		}
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// если получатель онлайн — доставляем по WS, иначе увидит при следующей загрузке
	h.hub.Deliver(msg)

	c.JSON(201, msg)
}

// GetHistory godoc
// @Summary      История сообщений с пользователем
// @Tags         messages
// @Produce      json
// @Param        id    path int false "ID собеседника"
// @Param        page  query int false "Страница" default(1)
// @Param        limit query int false "Лимит" default(50)
// @Success      200 {array} models.Message
// @Security     BearerAuth
// @Router       /messages/{id} [get]
func (h *MessageHandler) GetHistory(c *gin.Context) {
	userID := c.GetInt("user_id")
	partnerID, ok := parseIDParam(c, "id", "invalid id")
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}

	// помечаем прочитанными сообщения от партнёра
	if err := h.svc.MarkRead(userID, partnerID); err != nil {
		logger.Log.Warn().Err(err).Int("user_id", userID).Int("partner_id", partnerID).Msg("failed to mark messages as read")
	}

	msgs, err := h.svc.GetHistory(userID, partnerID, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	if msgs == nil {
		msgs = []models.Message{}
	}
	c.JSON(200, msgs)
}

// GetConversations godoc
// @Summary      Список переписок (чатов)
// @Tags         messages
// @Produce      json
// @Success      200 {array} models.ConversationPreview
// @Security     BearerAuth
// @Router       /messages [get]
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetInt("user_id")
	convs, err := h.svc.GetConversations(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	if convs == nil {
		convs = []models.ConversationPreview{}
	}
	c.JSON(200, convs)
}
