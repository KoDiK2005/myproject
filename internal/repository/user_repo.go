package repository

import (
	"context"
	"database/sql"
	"myproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRowxContext(ctx, query, user.Name, user.Email, user.PasswordHash).Scan(&user.ID)
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT id, name, email, avatar FROM users WHERE id = $1", id)
	return &user, err
}

func (r *UserRepo) GetAll(limit, offset int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var users []models.User
	err := r.db.SelectContext(ctx, &users, "SELECT id, name, email, avatar FROM users LIMIT $1 OFFSET $2", limit, offset)
	return users, err
}

func (r *UserRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepo) Update(id int, name, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var user models.User
	query := "UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING id, name, email, avatar"
	err := r.db.QueryRowxContext(ctx, query, name, email, id).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar)
	return &user, err
}

func (r *UserRepo) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var count int
	err := r.db.QueryRowxContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (r *UserRepo) UpdateAvatar(id int, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, "UPDATE users SET avatar = $1 WHERE id = $2", path, id)
	return err
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT id, name, email, password_hash FROM users WHERE email = $1", email)
	return &user, err
}
