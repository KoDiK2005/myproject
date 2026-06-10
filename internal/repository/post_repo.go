package repository

import (
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(post *models.Post) error {
	query := "INSERT INTO posts (title, body, user_id) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRowx(query, post.Title, post.Body, post.UserID).Scan(&post.ID)
}

func (r *PostRepo) GetByID(id int) (*models.Post, error) {
	var post models.Post
	err := r.db.Get(&post, "SELECT id, title, body, user_id FROM posts WHERE id = $1", id)
	return &post, err
}

func (r *PostRepo) GetAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Select(&posts, "SELECT id, title, body, user_id FROM posts LIMIT $1 OFFSET $2", limit, offset)
	return posts, err
}

func (r *PostRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}
