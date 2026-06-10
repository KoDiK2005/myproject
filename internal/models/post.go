package models

type Post struct {
	ID     int    `db:"id"      json:"id"`
	Title  string `db:"title"   json:"title"`
	Body   string `db:"body"    json:"body"`
	UserID int    `db:"user_id" json:"user_id"`
}

type CreatePostInput struct {
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body"  binding:"required"`
	UserID int    `json:"user_id" binding:"required"`
}
