package repository

import (
	"context"
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type FriendshipRepo struct {
	db *sqlx.DB
}

func NewFriendshipRepo(db *sqlx.DB) *FriendshipRepo {
	return &FriendshipRepo{db: db}
}

// SendRequest создаёт заявку в дружбы со статусом pending
func (r *FriendshipRepo) SendRequest(requesterID, addresseeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO friendships (requester_id, addressee_id, status)
		 VALUES ($1, $2, 'pending')
		 ON CONFLICT (requester_id, addressee_id) DO NOTHING`,
		requesterID, addresseeID,
	)
	return err
}

// Accept принимает заявку — addresseeID принимает заявку от requesterID
func (r *FriendshipRepo) Accept(requesterID, addresseeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE friendships SET status = 'accepted'
		 WHERE requester_id = $1 AND addressee_id = $2 AND status = 'pending'`,
		requesterID, addresseeID,
	)
	return err
}

// Reject отклоняет заявку — просто удаляем запись
func (r *FriendshipRepo) Reject(requesterID, addresseeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM friendships
		 WHERE requester_id = $1 AND addressee_id = $2 AND status = 'pending'`,
		requesterID, addresseeID,
	)
	return err
}

// Remove удаляет дружбу в любую сторону
func (r *FriendshipRepo) Remove(userID, friendID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM friendships
		 WHERE (requester_id = $1 AND addressee_id = $2)
		    OR (requester_id = $2 AND addressee_id = $1)`,
		userID, friendID,
	)
	return err
}

// GetFriends возвращает список принятых друзей юзера
func (r *FriendshipRepo) GetFriends(userID int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.id, u.name, u.email, u.avatar
		 FROM users u
		 JOIN friendships f ON (
		     (f.requester_id = $1 AND f.addressee_id = u.id)
		  OR (f.addressee_id = $1 AND f.requester_id = u.id)
		 )
		 WHERE f.status = 'accepted'
		 ORDER BY u.name`,
		userID,
	)
	return users, err
}

// GetIncomingRequests возвращает заявки, ожидающие ответа от userID
func (r *FriendshipRepo) GetIncomingRequests(userID int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.id, u.name, u.email, u.avatar
		 FROM users u
		 JOIN friendships f ON f.requester_id = u.id
		 WHERE f.addressee_id = $1 AND f.status = 'pending'
		 ORDER BY f.created_at DESC`,
		userID,
	)
	return users, err
}

// GetOutgoingRequests возвращает заявки которые userID отправил сам (ждут ответа)
func (r *FriendshipRepo) GetOutgoingRequests(userID int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		`SELECT u.id, u.name, u.email, u.avatar
		 FROM users u
		 JOIN friendships f ON f.addressee_id = u.id
		 WHERE f.requester_id = $1 AND f.status = 'pending'
		 ORDER BY f.created_at DESC`,
		userID,
	)
	return users, err
}

// GetStatus возвращает статус дружбы между двумя юзерами
func (r *FriendshipRepo) GetStatus(userID, otherID int) (*models.FriendshipStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var f models.Friendship
	err := r.db.GetContext(ctx, &f,
		`SELECT id, requester_id, addressee_id, status
		 FROM friendships
		 WHERE (requester_id = $1 AND addressee_id = $2)
		    OR (requester_id = $2 AND addressee_id = $1)`,
		userID, otherID,
	)
	if err != nil {
		// нет записи — не друзья и нет заявок
		return &models.FriendshipStatus{Status: "none"}, nil
	}

	status := f.Status
	if status == "pending" {
		if f.RequesterID == userID {
			status = "pending_sent"
		} else {
			status = "pending_received"
		}
	}
	return &models.FriendshipStatus{Status: status, FriendshipID: f.ID}, nil
}
