// ws/hub.go — WebSocket hub для личных сообщений.
// Одна горутина на hub, всё управление через каналы — никаких мьютексов.
package ws

import (
	"encoding/json"
	"myproject/internal/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 5120 // 5KB
)

// Client — один подключённый юзер
type Client struct {
	userID int
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
}

// Hub — хранит всех подключённых, раздаёт сообщения нужному получателю
type Hub struct {
	mu      sync.RWMutex
	clients map[int]*Client // userID → Client

	broadcast  chan *models.Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan *models.Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run — главная горутина хаба, запускаем один раз при старте сервера
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c.userID] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[c.userID]; ok && existing == c {
				delete(h.clients, c.userID)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			data, err := json.Marshal(models.WSMessage{Type: "message", Payload: msg})
			if err != nil {
				continue
			}
			// отправляем получателю если онлайн
			h.mu.RLock()
			receiver, online := h.clients[msg.ReceiverID]
			h.mu.RUnlock()
			if online {
				select {
				case receiver.send <- data:
				default:
					// буфер переполнен — отключаем
					h.mu.Lock()
					delete(h.clients, msg.ReceiverID)
					close(receiver.send)
					h.mu.Unlock()
				}
			}
		}
	}
}

// Deliver — отправить сообщение нужному юзеру если онлайн
func (h *Hub) Deliver(msg *models.Message) {
	h.broadcast <- msg
}

// writePump — пишет сообщения в WebSocket соединение
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump — читает входящие (держит соединение живым через pong)
// Реальная отправка сообщений — через REST, WS только для получения
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// NewClient — создаём клиента и запускаем горутины
func NewClient(hub *Hub, userID int, conn *websocket.Conn) {
	c := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
	hub.register <- c
	go c.writePump()
	go c.readPump()
}
