package models

import "time"

type Block struct {
	ID        int       `db:"id"         json:"id"`
	BlockerID int       `db:"blocker_id" json:"blocker_id"`
	BlockedID int       `db:"blocked_id" json:"blocked_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
