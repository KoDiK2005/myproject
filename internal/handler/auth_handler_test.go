package handler_test

import (
	"bytes"
	"encoding/json"
	"myproject/internal/handler"
	"myproject/internal/models"
	"myproject/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo — заглушка юзер-репозитория
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
	return nil, nil
}

func (m *mockUserRepo) GetAll(limit, offset int) ([]models.User, error) {
	return m.users, nil
}

func (m *mockUserRepo) Count() (int, error) {
	return len(m.users), nil
}

func (m *mockUserRepo) Delete(id int) error { return nil }

func (m *mockUserRepo) Update(id int, name, email string) (*models.User, error) {
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) UpdateAvatar(id int, path string) error { return nil }

func (m *mockUserRepo) Search(query string, limit, offset int) ([]models.User, int, error) {
	return m.users, len(m.users), nil
}

// mockRefreshTokenRepo — минимальная реализация
type mockRefreshTokenRepo struct {
	tokens []models.RefreshToken
}

func (m *mockRefreshTokenRepo) Create(token *models.RefreshToken) error {
	token.ID = len(m.tokens) + 1
	m.tokens = append(m.tokens, *token)
	return nil
}

func (m *mockRefreshTokenRepo) GetByToken(tok string) (*models.RefreshToken, error) {
	for _, t := range m.tokens {
		if t.Token == tok {
			return &t, nil
		}
	}
	return nil, nil
}

func (m *mockRefreshTokenRepo) Revoke(tok string) error {
	for i, t := range m.tokens {
		if t.Token == tok {
			m.tokens[i].Revoked = true
		}
	}
	return nil
}

func (m *mockRefreshTokenRepo) RevokeAllForUser(userID int) error {
	for i, t := range m.tokens {
		if t.UserID == userID {
			m.tokens[i].Revoked = true
		}
	}
	return nil
}

// setupAuthRouter собирает роутер для тестов логина
func setupAuthRouter(userRepo service.UserRepository, rtRepo service.RefreshTokenRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	userSvc := service.NewUserService(userRepo)
	rtSvc := service.NewRefreshTokenService(rtRepo)
	h := handler.NewAuthHandler(userSvc, rtSvc, "test-secret")
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout-all", handler.AuthMiddleware("test-secret"), h.LogoutAll)
	return r
}

// делаем валидный JWT так же, как AuthHandler.Login
func makeTestAccessToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	str, _ := token.SignedString([]byte("test-secret"))
	return str
}

// хелпер — создаём юзера с захэшированным паролем
func newUserWithPassword(id int, email, password string) models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return models.User{
		ID:           id,
		Name:         "Testuser",
		Email:        email,
		PasswordHash: string(hash),
	}
}

func TestLogin_Success(t *testing.T) {
	user := newUserWithPassword(1, "test@example.com", "secret123")
	userRepo := &mockUserRepo{users: []models.User{user}}
	rtRepo := &mockRefreshTokenRepo{}

	r := setupAuthRouter(userRepo, rtRepo)
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "secret123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("ожидали 200, получили %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["access_token"] == "" {
		t.Error("access_token пустой — что-то пошло не так")
	}
	if resp["refresh_token"] == "" {
		t.Error("refresh_token пустой — что-то пошло не так")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := newUserWithPassword(1, "test@example.com", "secret123")
	userRepo := &mockUserRepo{users: []models.User{user}}
	rtRepo := &mockRefreshTokenRepo{}

	r := setupAuthRouter(userRepo, rtRepo)
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "неверный"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("ожидали 401, получили %d", w.Code)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	r := setupAuthRouter(&mockUserRepo{}, &mockRefreshTokenRepo{})
	body, _ := json.Marshal(map[string]string{"email": "noone@example.com", "password": "123456"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("ожидали 401 для несуществующего юзера, получили %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	r := setupAuthRouter(&mockUserRepo{}, &mockRefreshTokenRepo{})
	body, _ := json.Marshal(map[string]string{"email": "test@example.com"}) // нет password
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("ожидали 400 без пароля, получили %d", w.Code)
	}
}

func TestLogoutAll_RevokesAllUserTokens(t *testing.T) {
	rtRepo := &mockRefreshTokenRepo{tokens: []models.RefreshToken{
		{ID: 1, UserID: 1, Token: "tok-a"},
		{ID: 2, UserID: 1, Token: "tok-b"},
		{ID: 3, UserID: 2, Token: "tok-other-user"}, // не должен затронуться
	}}
	r := setupAuthRouter(&mockUserRepo{}, rtRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout-all", nil)
	req.Header.Set("Authorization", "Bearer "+makeTestAccessToken(1))
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Fatalf("ожидали 204, получили %d: %s", w.Code, w.Body.String())
	}
	for _, tok := range rtRepo.tokens {
		if tok.UserID == 1 && !tok.Revoked {
			t.Errorf("токен %q юзера 1 должен быть отозван", tok.Token)
		}
		if tok.UserID == 2 && tok.Revoked {
			t.Error("токен другого юзера не должен быть отозван")
		}
	}
}

func TestLogoutAll_RequiresAuth(t *testing.T) {
	r := setupAuthRouter(&mockUserRepo{}, &mockRefreshTokenRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout-all", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("ожидали 401 без токена, получили %d", w.Code)
	}
}

// маленький тест, что bcrypt MinCost работает быстро (не тормозим CI)
func TestBcryptMinCostFast(t *testing.T) {
	start := time.Now()
	newUserWithPassword(1, "x@x.com", "pass")
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("bcrypt.MinCost слишком медленный в тестах — %v", time.Since(start))
	}
}
