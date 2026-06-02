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
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRowx(query, user.Name, user.Email).Scan(&user.ID)
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, "SELECT id, name, email FROM users")
	return users, err
}
