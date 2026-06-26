package repository

import (
	"context"
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type BlockRepo struct {
	db *sqlx.DB
}

func NewBlockRepo(db *sqlx.DB) *BlockRepo {
	return &BlockRepo{db: db}
}

func (r *BlockRepo) Block(blockerID, blockedID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO blocks (blocker_id, blocked_id)
		 VALUES ($1, $2)
		 ON CONFLICT (blocker_id, blocked_id) DO NOTHING`,
		blockerID, blockedID,
	)
	return err
}

func (r *BlockRepo) Unblock(blockerID, blockedID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2`,
		blockerID, blockedID,
	)
	return err
}

// IsBlocked — true если есть блокировка в любую сторону между двумя юзерами
func (r *BlockRepo) IsBlocked(userA, userB int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		`SELECT EXISTS(
		   SELECT 1 FROM blocks
		   WHERE (blocker_id = $1 AND blocked_id = $2)
		      OR (blocker_id = $2 AND blocked_id = $1)
		 )`,
		userA, userB).Scan(&exists)
	return exists, err
}

// GetBlockedUsers — те, кого blockerID заблокировал сам
func (r *BlockRepo) GetBlockedUsers(blockerID int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.id, u.name, u.email, u.avatar
		 FROM users u
		 JOIN blocks b ON b.blocked_id = u.id
		 WHERE b.blocker_id = $1
		 ORDER BY b.created_at DESC`,
		blockerID,
	)
	return users, err
}
