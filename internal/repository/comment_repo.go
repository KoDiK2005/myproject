package repository

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `WITH inserted AS (
	            INSERT INTO comments (post_id, user_id, body)
	            VALUES ($1, $2, $3)
	            RETURNING id, created_at
	          )
	          SELECT inserted.id, inserted.created_at, users.name AS author_name
	          FROM inserted JOIN users ON users.id = $2`
	return r.db.QueryRowxContext(ctx, query, c.PostID, c.UserID, c.Body).
		Scan(&c.ID, &c.CreatedAt, &c.AuthorName)
}

func (r *CommentRepo) GetByPostID(postID, limit, offset int) ([]models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var comments []models.Comment
	err := r.db.SelectContext(ctx, &comments,
		`SELECT comments.id, comments.post_id, comments.user_id, comments.body, comments.created_at,
		        users.name AS author_name
		 FROM comments JOIN users ON users.id = comments.user_id
		 WHERE comments.post_id = $1
		 ORDER BY comments.created_at ASC
		 LIMIT $2 OFFSET $3`,
		postID, limit, offset)
	return comments, err
}

func (r *CommentRepo) CountByPostID(postID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var count int
	err := r.db.QueryRowxContext(ctx,
		`SELECT COUNT(*) FROM comments WHERE post_id = $1`, postID).Scan(&count)
	return count, err
}

func (r *CommentRepo) GetByID(id int) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var c models.Comment
	err := r.db.GetContext(ctx, &c,
		`SELECT id, post_id, user_id, body, created_at FROM comments WHERE id = $1`, id)
	return &c, err
}

func (r *CommentRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, "DELETE FROM comments WHERE id = $1", id)
	return err
}
