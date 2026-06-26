//go:build integration

// Запускать так: go test -tags integration ./internal/repository/... -v
// Требует работающий PostgreSQL через env-переменную DATABASE_URL

package repository

import (
	"testing"
)

// TestFriendshipRepo_MutualRequestAutoAccepts проверяет фикс бага: если A отправляет
// заявку B, а затем B отправляет встречную заявку A до того как A её принял,
// должна получиться ОДНА accepted-запись, а не две pending в разные стороны
// (что раньше приводило к дублям в GetFriends).
func TestFriendshipRepo_MutualRequestAutoAccepts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userA := createTestUser(t, db, "friendship-a@integration.test")
	userB := createTestUser(t, db, "friendship-b@integration.test")
	defer func() {
		db.Exec("DELETE FROM friendships WHERE requester_id IN ($1, $2) OR addressee_id IN ($1, $2)", userA, userB)
		db.Exec("DELETE FROM users WHERE id IN ($1, $2)", userA, userB)
	}()

	repo := NewFriendshipRepo(db)

	autoAccepted, err := repo.SendRequest(userA, userB)
	if err != nil {
		t.Fatalf("SendRequest(A,B) failed: %v", err)
	}
	if autoAccepted {
		t.Fatal("первая заявка не должна авто-приниматься")
	}

	autoAccepted, err = repo.SendRequest(userB, userA)
	if err != nil {
		t.Fatalf("SendRequest(B,A) failed: %v", err)
	}
	if !autoAccepted {
		t.Fatal("встречная заявка должна авто-принять существующую pending-запись")
	}

	var count int
	if err := db.Get(&count,
		`SELECT COUNT(*) FROM friendships
		 WHERE (requester_id = $1 AND addressee_id = $2) OR (requester_id = $2 AND addressee_id = $1)`,
		userA, userB,
	); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("ожидали ровно одну запись дружбы между A и B, получили %d", count)
	}

	var status string
	if err := db.Get(&status, `SELECT status FROM friendships WHERE requester_id = $1 AND addressee_id = $2`, userA, userB); err != nil {
		t.Fatalf("status query failed: %v", err)
	}
	if status != "accepted" {
		t.Errorf("ожидали status=accepted, получили %q", status)
	}

	friends, err := repo.GetFriends(userA)
	if err != nil {
		t.Fatalf("GetFriends failed: %v", err)
	}
	seen := 0
	for _, f := range friends {
		if f.ID == userB {
			seen++
		}
	}
	if seen != 1 {
		t.Errorf("ожидали что B встречается в списке друзей A ровно раз, получили %d", seen)
	}
}
