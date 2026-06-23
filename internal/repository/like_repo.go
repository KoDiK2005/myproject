package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type LikeRepo struct {
	db *sqlx.DB
}

func NewLikeRepo(db *sqlx.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

// Like добавляет лайк; если уже есть — молча игнорирует (ON CONFLICT DO NOTHING)
func (r *LikeRepo) Like(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID,
	)
	return err
}

func (r *LikeRepo) Unlike(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM likes WHERE user_id = $1 AND post_id = $2`,
		userID, postID,
	)
	return err
}

func (r *LikeRepo) Count(postID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var count int
	err := r.db.QueryRowxContext(ctx, `SELECT COUNT(*) FROM likes WHERE post_id = $1`, postID).Scan(&count)
	return count, err
}

// IsLiked — лайкнул ли userID пост postID
func (r *LikeRepo) IsLiked(userID, postID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`,
		userID, postID).Scan(&exists)
	return exists, err
}
