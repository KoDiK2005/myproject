package repository

import "github.com/jmoiron/sqlx"

type LikeRepo struct {
	db *sqlx.DB
}

func NewLikeRepo(db *sqlx.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

// Like добавляет лайк; если уже есть — молча игнорирует (ON CONFLICT DO NOTHING)
func (r *LikeRepo) Like(userID, postID int) error {
	_, err := r.db.Exec(
		`INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID,
	)
	return err
}

func (r *LikeRepo) Unlike(userID, postID int) error {
	_, err := r.db.Exec(
		`DELETE FROM likes WHERE user_id = $1 AND post_id = $2`,
		userID, postID,
	)
	return err
}

func (r *LikeRepo) Count(postID int) (int, error) {
	var count int
	err := r.db.QueryRowx(`SELECT COUNT(*) FROM likes WHERE post_id = $1`, postID).Scan(&count)
	return count, err
}
