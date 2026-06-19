package repository

import (
	"context"
	"myproject/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MessageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

// Save — сохраняем сообщение в БД
func (r *MessageRepo) Save(msg *models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO messages (sender_id, receiver_id, content)
		 VALUES ($1, $2, $3) RETURNING id, created_at`,
		msg.SenderID, msg.ReceiverID, msg.Content,
	).Scan(&msg.ID, &msg.CreatedAt)
}

// GetHistory — история сообщений между двумя юзерами, от старых к новым
func (r *MessageRepo) GetHistory(userA, userB, limit, offset int) ([]models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var msgs []models.Message
	err := r.db.SelectContext(ctx, &msgs,
		`SELECT id, sender_id, receiver_id, content, read_at, created_at
		 FROM messages
		 WHERE (sender_id = $1 AND receiver_id = $2)
		    OR (sender_id = $2 AND receiver_id = $1)
		 ORDER BY created_at ASC
		 LIMIT $3 OFFSET $4`,
		userA, userB, limit, offset)
	return msgs, err
}

// GetConversations — список всех чатов с последним сообщением
func (r *MessageRepo) GetConversations(userID int) ([]models.ConversationPreview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var list []models.ConversationPreview
	err := r.db.SelectContext(ctx, &list,
		`SELECT
		   partner.id       AS user_id,
		   partner.name     AS user_name,
		   partner.avatar   AS user_avatar,
		   last_msg.content AS last_message,
		   last_msg.created_at AS last_at,
		   COALESCE(unread.cnt, 0) AS unread
		 FROM (
		   -- все собеседники текущего юзера
		   SELECT DISTINCT
		     CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END AS partner_id
		   FROM messages
		   WHERE sender_id = $1 OR receiver_id = $1
		 ) AS partners
		 JOIN users AS partner ON partner.id = partners.partner_id
		 -- последнее сообщение с каждым
		 JOIN LATERAL (
		   SELECT content, created_at FROM messages
		   WHERE (sender_id = $1 AND receiver_id = partners.partner_id)
		      OR (sender_id = partners.partner_id AND receiver_id = $1)
		   ORDER BY created_at DESC
		   LIMIT 1
		 ) AS last_msg ON true
		 -- непрочитанные от партнёра
		 LEFT JOIN LATERAL (
		   SELECT COUNT(*) AS cnt FROM messages
		   WHERE sender_id = partners.partner_id
		     AND receiver_id = $1
		     AND read_at IS NULL
		 ) AS unread ON true
		 ORDER BY last_msg.created_at DESC`,
		userID)
	return list, err
}

// MarkRead — помечаем все сообщения от sender к receiver как прочитанные
func (r *MessageRepo) MarkRead(receiverID, senderID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE messages SET read_at = $1
		 WHERE receiver_id = $2 AND sender_id = $3 AND read_at IS NULL`,
		now, receiverID, senderID)
	return err
}

// IsFriend — проверка что юзеры друзья (только друзья могут писать)
func (r *MessageRepo) IsFriend(userA, userB int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		`SELECT EXISTS(
		   SELECT 1 FROM friendships
		   WHERE ((requester_id = $1 AND addressee_id = $2)
		       OR (requester_id = $2 AND addressee_id = $1))
		   AND status = 'accepted'
		 )`,
		userA, userB).Scan(&exists)
	return exists, err
}
