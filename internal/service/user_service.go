package service

import (
	"errors"
	"myproject/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	Delete(id int) error
	Update(id int, name, email string) (*models.User, error) // имя и мейл — не весь юзер
	GetByEmail(email string) (*models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(input models.CreateUserInput) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}
	user.PasswordHash = string(hash)

	err = s.repo.Create(user)
	return user, err
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ListUsers(page, limit int) ([]models.User, error) {
	offset := (page - 1) * limit
	return s.repo.GetAll(limit, offset)
}
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) UpdateUser(id int, input models.UpdateUserInput) (*models.User, error) {
	return s.repo.Update(id, input.Name, input.Email)
}

func (s *UserService) Login(input models.LoginInput) (*models.User, error) {
	user, err := s.repo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
