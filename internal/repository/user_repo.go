package repository

import (
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
	query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRowx(query, user.Name, user.Email, user.PasswordHash).Scan(&user.ID)
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT id, name, email FROM users WHERE id = $1", id)
	return &user, err
}

func (r *UserRepo) GetAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, "SELECT id, name, email FROM users LIMIT $1 OFFSET $2", limit, offset)
	return users, err
}

func (r *UserRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *UserRepo) Update(id int, name, email string) (*models.User, error) {
	var user models.User
	query := "UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING id, name, email"
	err := r.db.QueryRowx(query, name, email, id).Scan(&user.ID, &user.Name, &user.Email)
	return &user, err
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT id, name, email, password_hash FROM users WHERE email = $1", email)
	return &user, err
}
