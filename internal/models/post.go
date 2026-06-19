package models

type Post struct {
	ID         int    `db:"id"         json:"id"`
	Title      string `db:"title"      json:"title"`
	Body       string `db:"body"       json:"body"`
	UserID     int    `db:"user_id"    json:"user_id"`
	Visibility string `db:"visibility" json:"visibility"` // public | friends
}

type CreatePostInput struct {
	Title      string `json:"title"      binding:"required,min=3,max=200"`
	Body       string `json:"body"       binding:"required,min=10,max=10000"`
	Visibility string `json:"visibility"` // public | friends, default public
}

type UpdatePostInput struct {
	Title      string `json:"title"      binding:"required,min=3,max=200"`
	Body       string `json:"body"       binding:"required,min=10,max=10000"`
	Visibility string `json:"visibility"`
}
