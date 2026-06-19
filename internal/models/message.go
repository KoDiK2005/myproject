package models

import "time"

type Message struct {
	ID         int        `db:"id"          json:"id"`
	SenderID   int        `db:"sender_id"   json:"sender_id"`
	ReceiverID int        `db:"receiver_id" json:"receiver_id"`
	Content    string     `db:"content"     json:"content"`
	ReadAt     *time.Time `db:"read_at"     json:"read_at,omitempty"`
	CreatedAt  time.Time  `db:"created_at"  json:"created_at"`
}

// ConversationPreview — краткая инфа о чате для списка переписок
type ConversationPreview struct {
	UserID      int        `db:"user_id"      json:"user_id"`
	UserName    string     `db:"user_name"    json:"user_name"`
	UserAvatar  *string    `db:"user_avatar"  json:"user_avatar,omitempty"`
	LastMessage string     `db:"last_message" json:"last_message"`
	LastAt      time.Time  `db:"last_at"      json:"last_at"`
	Unread      int        `db:"unread"       json:"unread"`
}

// WSMessage — то что летит по WebSocket
type WSMessage struct {
	Type    string   `json:"type"` // "message" | "read"
	Payload *Message `json:"payload,omitempty"`
}
