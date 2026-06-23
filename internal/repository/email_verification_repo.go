package repository

import (
	"context"
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type EmailVerificationRepo struct {
	db *sqlx.DB
}

func NewEmailVerificationRepo(db *sqlx.DB) *EmailVerificationRepo {
	return &EmailVerificationRepo{db: db}
}

func (r *EmailVerificationRepo) Create(token *models.EmailVerificationToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO email_verification_tokens (user_id, token, expires_at)
		 VALUES ($1, $2, $3) RETURNING id`,
		token.UserID, token.Token, token.ExpiresAt,
	).Scan(&token.ID)
}

func (r *EmailVerificationRepo) GetByToken(token string) (*models.EmailVerificationToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var t models.EmailVerificationToken
	err := r.db.GetContext(ctx, &t,
		"SELECT id, user_id, token, expires_at, used FROM email_verification_tokens WHERE token = $1", token)
	return &t, err
}

func (r *EmailVerificationRepo) MarkUsed(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, "UPDATE email_verification_tokens SET used = true WHERE token = $1", token)
	return err
}
