package models

import "time"

type Friendship struct {
	ID          int       `db:"id"           json:"id"`
	RequesterID int       `db:"requester_id" json:"requester_id"`
	AddresseeID int       `db:"addressee_id" json:"addressee_id"`
	Status      string    `db:"status"       json:"status"` // pending | accepted | rejected
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

// FriendshipStatus — то что фронт получает для конкретной пары юзеров
type FriendshipStatus struct {
	Status      string `json:"status"`       // none | pending_sent | pending_received | accepted
	FriendshipID int  `json:"friendship_id"` // 0 если нет
}
