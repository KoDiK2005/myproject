package models

import "time"

// Notification — type: friend_request | friend_accept | like | comment
type Notification struct {
	ID        int       `db:"id"         json:"id"`
	UserID    int       `db:"user_id"    json:"user_id"`
	ActorID   int       `db:"actor_id"   json:"actor_id"`
	ActorName string    `db:"actor_name" json:"actor_name"`
	Type      string    `db:"type"       json:"type"`
	PostID    *int      `db:"post_id"    json:"post_id,omitempty"`
	Read      bool      `db:"read"       json:"read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// WSNotification — то что летит по WebSocket при новом уведомлении
type WSNotification struct {
	Type    string        `json:"type"` // всегда "notification"
	Payload *Notification `json:"payload"`
}
