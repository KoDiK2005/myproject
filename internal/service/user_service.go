package service

import (
	"errors"
	"myproject/internal/models"
	"myproject/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(input models.CreateUserInput) (*models.User, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err := s.repo.Create(user)
	return user, err
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.repo.GetAll()
}
