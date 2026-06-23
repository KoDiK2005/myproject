package models

import "time"

type Comment struct {
	ID         int       `db:"id"          json:"id"`
	PostID     int       `db:"post_id"     json:"post_id"`
	UserID     int       `db:"user_id"     json:"user_id"`
	Body       string    `db:"body"        json:"body"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	AuthorName string    `db:"author_name" json:"author_name"`
}

type CreateCommentInput struct {
	Body string `json:"body" binding:"required,min=1,max=1000"`
}
