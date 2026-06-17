package repository

import (
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(c *models.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, body)
	          VALUES ($1, $2, $3)
	          RETURNING id, created_at`
	return r.db.QueryRowx(query, c.PostID, c.UserID, c.Body).Scan(&c.ID, &c.CreatedAt)
}

func (r *CommentRepo) GetByPostID(postID int) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Select(&comments,
		`SELECT id, post_id, user_id, body, created_at
		 FROM comments WHERE post_id = $1 ORDER BY created_at ASC`,
		postID)
	return comments, err
}

func (r *CommentRepo) GetByID(id int) (*models.Comment, error) {
	var c models.Comment
	err := r.db.Get(&c,
		`SELECT id, post_id, user_id, body, created_at FROM comments WHERE id = $1`, id)
	return &c, err
}

func (r *CommentRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE id = $1", id)
	return err
}
