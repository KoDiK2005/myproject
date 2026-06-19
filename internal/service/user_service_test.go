package service

import (
	"errors"
	"myproject/internal/models"
	"testing"
)

type mockUserRepo struct {
	users []models.User
}

func (m *mockUserRepo) Create(user *models.User) error {
	user.ID = len(m.users) + 1
	m.users = append(m.users, *user)
	return nil
}

func (m *mockUserRepo) GetByID(id int) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetAll(limit, offset int) ([]models.User, error) {
	return m.users, nil
}

func (m *mockUserRepo) Count() (int, error) {
	return len(m.users), nil
}

func (m *mockUserRepo) Delete(id int) error {
	return nil
}

func (m *mockUserRepo) Update(id int, name, email string) (*models.User, error) {
	for i, u := range m.users {
		if u.ID == id {
			m.users[i].Name = name
			m.users[i].Email = email
			return &m.users[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Search(query string, limit, offset int) ([]models.User, int, error) {
	return m.users, len(m.users), nil
}

func (m *mockUserRepo) UpdateAvatar(id int, path string) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users[i].Avatar = path
			return nil
		}
	}
	return errors.New("not found")
}

func TestCreateUser(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	input := models.CreateUserInput{Name: "Mark", Email: "mark@example.com", Password: "secret123"}
	user, err := svc.CreateUser(input)

	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if user.Name != "Mark" {
		t.Errorf("ожидали name='Mark', получили '%s'", user.Name)
	}
	if user.Email != "mark@example.com" {
		t.Errorf("ожидали email='mark@example.com', получили '%s'", user.Email)
	}
}

func TestListUsers(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	svc.CreateUser(models.CreateUserInput{Name: "Mark", Email: "mark@example.com", Password: "123"})
	svc.CreateUser(models.CreateUserInput{Name: "Anna", Email: "anna@example.com", Password: "123"})

	resp, err := svc.ListUsers(1, 10)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	users := resp.Data.([]models.User)
	if len(users) != 2 {
		t.Errorf("ожидали 2 юзера, получили %d", len(users))
	}
	if resp.Total != 2 {
		t.Errorf("ожидали total=2, получили %d", resp.Total)
	}
}
