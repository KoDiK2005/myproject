//go:build integration

// Запускать так: go test -tags integration ./internal/repository/... -v
// Требует работающий PostgreSQL через env-переменную DATABASE_URL

package repository

import (
	"myproject/internal/models"
	"testing"

	"github.com/jmoiron/sqlx"
)

func createTestUser(t *testing.T, db *sqlx.DB, email string) int {
	t.Helper()
	var id int
	err := db.QueryRowx(
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		email, email, "hash",
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test user %s: %v", email, err)
	}
	return id
}

func containsPostID(posts []models.Post, id int) bool {
	for _, p := range posts {
		if p.ID == id {
			return true
		}
	}
	return false
}

// TestPostRepo_BlockedUsersExcludedEverywhere проверяет, что блокировка скрывает посты
// от/для заблокированного во всех путях чтения — фид, поиск, список постов профиля.
// Раньше блокировка влияла только на прямой запрос поста по ID.
func TestPostRepo_BlockedUsersExcludedEverywhere(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	author := createTestUser(t, db, "post-author@integration.test")
	blocker := createTestUser(t, db, "post-blocker@integration.test")
	bystander := createTestUser(t, db, "post-bystander@integration.test")
	defer func() {
		db.Exec("DELETE FROM users WHERE id IN ($1, $2, $3)", author, blocker, bystander)
	}()

	repo := NewPostRepo(db)

	postID := 0
	err := db.QueryRowx(
		"INSERT INTO posts (title, body, user_id, visibility) VALUES ($1, $2, $3, 'public') RETURNING id",
		"unique-searchable-title", "body", author,
	).Scan(&postID)
	if err != nil {
		t.Fatalf("failed to insert post: %v", err)
	}
	defer db.Exec("DELETE FROM posts WHERE id = $1", postID)

	if _, err := db.Exec("INSERT INTO blocks (blocker_id, blocked_id) VALUES ($1, $2)", blocker, author); err != nil {
		t.Fatalf("failed to insert block: %v", err)
	}
	defer db.Exec("DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2", blocker, author)

	t.Run("feed excludes blocked author's post", func(t *testing.T) {
		posts, _, err := repo.GetFeedWithCount(blocker, 50, 0)
		if err != nil {
			t.Fatalf("GetFeedWithCount failed: %v", err)
		}
		if containsPostID(posts, postID) {
			t.Error("ожидали что пост заблокированного автора скрыт из фида")
		}

		posts, _, err = repo.GetFeedWithCount(bystander, 50, 0)
		if err != nil {
			t.Fatalf("GetFeedWithCount failed: %v", err)
		}
		if !containsPostID(posts, postID) {
			t.Error("посторонний пользователь должен видеть публичный пост")
		}
	})

	t.Run("search excludes blocked author's post", func(t *testing.T) {
		posts, _, err := repo.SearchWithCount("unique-searchable-title", blocker, 50, 0)
		if err != nil {
			t.Fatalf("SearchWithCount failed: %v", err)
		}
		if containsPostID(posts, postID) {
			t.Error("ожидали что пост заблокированного автора скрыт из поиска")
		}

		posts, _, err = repo.SearchWithCount("unique-searchable-title", bystander, 50, 0)
		if err != nil {
			t.Fatalf("SearchWithCount failed: %v", err)
		}
		if !containsPostID(posts, postID) {
			t.Error("посторонний пользователь должен находить публичный пост в поиске")
		}
	})

	t.Run("profile post list excludes blocked author's post", func(t *testing.T) {
		posts, err := repo.GetByUserIDForViewer(author, blocker, 50, 0)
		if err != nil {
			t.Fatalf("GetByUserIDForViewer failed: %v", err)
		}
		if containsPostID(posts, postID) {
			t.Error("ожидали что заблокированный не видит пост в профиле автора")
		}

		posts, err = repo.GetByUserIDForViewer(author, bystander, 50, 0)
		if err != nil {
			t.Fatalf("GetByUserIDForViewer failed: %v", err)
		}
		if !containsPostID(posts, postID) {
			t.Error("посторонний должен видеть пост в профиле автора")
		}
	})
}
