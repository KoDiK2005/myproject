package service

import (
	"database/sql"
	"errors"
	"myproject/internal/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	Search(query string, limit, offset int) ([]models.User, int, error)
	Count() (int, error)
	Delete(id int) error
	Update(id int, name, email string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateAvatar(id int, path string) error
	MarkEmailVerified(userID int) error
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
	if err != nil {
		// postgres unique violation — email уже занят
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ListUsers(page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	users, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: users, Total: total, Page: page, Limit: limit}, nil
}
func (s *UserService) SearchUsers(query string, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	users, total, err := s.repo.Search(query, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: users, Total: total, Page: page, Limit: limit}, nil
}

func (s *UserService) DeleteUser(id int) error {
	err := s.repo.Delete(id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *UserService) UpdateUser(id int, input models.UpdateUserInput) (*models.User, error) {
	user, err := s.repo.Update(id, input.Name, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateAvatar(id int, path string) error {
	return s.repo.UpdateAvatar(id, path)
}

func (s *UserService) Login(input models.LoginInput) (*models.User, error) {
	user, err := s.repo.GetByEmail(input.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
