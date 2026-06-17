//go:build integration

// Запускать так: go test -tags integration ./internal/repository/... -v
// Требует работающий PostgreSQL через env-переменную DATABASE_URL
// Например: DATABASE_URL=postgres://postgres:password@localhost:5432/testdb?sslmode=disable

package repository

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL не задан, пропускаем интеграционный тест")
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("не удалось подключиться к БД: %v", err)
	}
	return db
}

func TestUserRepo_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// чистим таблицу перед тестом чтобы не было мусора
	db.Exec("DELETE FROM users WHERE email = $1", "integration@test.com")

	repo := NewUserRepo(db)

	user := &struct {
		ID           int
		Name         string
		Email        string
		PasswordHash string
	}{Name: "Integration", Email: "integration@test.com", PasswordHash: "hash"}

	// используем напрямую SQL, потому что Create принимает *models.User
	var id int
	err := db.QueryRowx(
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.PasswordHash,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	got, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Email != user.Email {
		t.Errorf("email mismatch: got %s, want %s", got.Email, user.Email)
	}

	// чистим за собой
	db.Exec("DELETE FROM users WHERE id = $1", id)
}

func TestUserRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec("DELETE FROM users WHERE email = $1", "del@test.com")

	repo := NewUserRepo(db)

	var id int
	db.QueryRowx(
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		"ToDelete", "del@test.com", "hash",
	).Scan(&id)

	if err := repo.Delete(id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.GetByID(id)
	if err == nil {
		t.Error("ожидали ошибку после удаления, но получили nil")
	}
}
