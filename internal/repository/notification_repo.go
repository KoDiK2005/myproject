package repository

import (
	"context"
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type NotificationRepo struct {
	db *sqlx.DB
}

func NewNotificationRepo(db *sqlx.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(n *models.Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO notifications (user_id, actor_id, type, post_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		n.UserID, n.ActorID, n.Type, n.PostID,
	).Scan(&n.ID, &n.CreatedAt)
}

func (r *NotificationRepo) GetByUserID(userID, limit, offset int) ([]models.Notification, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var list []models.Notification
	err := r.db.SelectContext(ctx, &list,
		`SELECT n.id, n.user_id, n.actor_id, u.name AS actor_name, n.type, n.post_id, n.read, n.created_at
		 FROM notifications n
		 JOIN users u ON u.id = n.actor_id
		 WHERE n.user_id = $1
		 ORDER BY n.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *NotificationRepo) MarkRead(userID, notificationID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2`,
		notificationID, userID,
	)
	return err
}

func (r *NotificationRepo) MarkAllRead(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`,
		userID,
	)
	return err
}

func (r *NotificationRepo) UnreadCount(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`,
		userID,
	)
	return count, err
}
