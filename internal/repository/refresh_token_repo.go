package repository

import (
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepo(db *sqlx.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(token *models.RefreshToken) error {
	_, err := r.db.Exec(
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt,
	)
	return err
}

func (r *RefreshTokenRepo) GetByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Get(&rt, "SELECT id, user_id, token, expires_at, revoked FROM refresh_tokens WHERE token = $1", token)
	return &rt, err
}

func (r *RefreshTokenRepo) Revoke(token string) error {
	_, err := r.db.Exec("UPDATE refresh_tokens SET revoked = true WHERE token = $1", token)
	return err
}
