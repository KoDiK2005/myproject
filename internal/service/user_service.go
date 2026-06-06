package service

import (
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

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) UpdateUser(id int, input models.CreateUserInput) (*models.User, error) {
	return s.repo.Update(id, input.Name, input.Email)
}
