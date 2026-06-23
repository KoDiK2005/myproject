package service

import (
	"errors"
	"myproject/internal/models"
	"testing"
	"time"
)

type mockEmailVerificationRepo struct {
	tokens []models.EmailVerificationToken
}

func (m *mockEmailVerificationRepo) Create(token *models.EmailVerificationToken) error {
	token.ID = len(m.tokens) + 1
	m.tokens = append(m.tokens, *token)
	return nil
}

func (m *mockEmailVerificationRepo) GetByToken(token string) (*models.EmailVerificationToken, error) {
	for _, t := range m.tokens {
		if t.Token == token {
			return &t, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockEmailVerificationRepo) MarkUsed(token string) error {
	for i, t := range m.tokens {
		if t.Token == token {
			m.tokens[i].Used = true
			return nil
		}
	}
	return errors.New("not found")
}

type mockEmailSender struct {
	sentTo      string
	sentSubject string
	sentBody    string
}

func (m *mockEmailSender) Send(to, subject, body string) error {
	m.sentTo, m.sentSubject, m.sentBody = to, subject, body
	return nil
}

func TestSendVerification_GeneratesTokenAndSendsEmail(t *testing.T) {
	repo := &mockEmailVerificationRepo{}
	userRepo := &mockUserRepo{users: []models.User{{ID: 1, Email: "test@example.com"}}}
	sender := &mockEmailSender{}
	svc := NewEmailVerificationService(repo, userRepo, sender, "http://localhost:5173")

	err := svc.SendVerification(1, "test@example.com")
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if len(repo.tokens) != 1 {
		t.Fatalf("ожидали 1 токен, получили %d", len(repo.tokens))
	}
	if sender.sentTo != "test@example.com" {
		t.Errorf("письмо отправлено не туда: %q", sender.sentTo)
	}
}

func TestVerifyToken_Success(t *testing.T) {
	repo := &mockEmailVerificationRepo{tokens: []models.EmailVerificationToken{
		{Token: "good-token", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
	}}
	userRepo := &mockUserRepo{users: []models.User{{ID: 1, Email: "test@example.com"}}}
	svc := NewEmailVerificationService(repo, userRepo, &mockEmailSender{}, "http://localhost:5173")

	if err := svc.VerifyToken("good-token"); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if !userRepo.users[0].EmailVerified {
		t.Error("email должен быть помечен подтверждённым")
	}
	if !repo.tokens[0].Used {
		t.Error("токен должен быть помечен использованным")
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	repo := &mockEmailVerificationRepo{tokens: []models.EmailVerificationToken{
		{Token: "old-token", UserID: 1, ExpiresAt: time.Now().Add(-time.Hour)},
	}}
	userRepo := &mockUserRepo{users: []models.User{{ID: 1}}}
	svc := NewEmailVerificationService(repo, userRepo, &mockEmailSender{}, "http://localhost:5173")

	if err := svc.VerifyToken("old-token"); !errors.Is(err, ErrTokenInvalid) {
		t.Errorf("ожидали ErrTokenInvalid для просроченного токена, получили: %v", err)
	}
}

func TestVerifyToken_AlreadyUsed(t *testing.T) {
	repo := &mockEmailVerificationRepo{tokens: []models.EmailVerificationToken{
		{Token: "used-token", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), Used: true},
	}}
	userRepo := &mockUserRepo{users: []models.User{{ID: 1}}}
	svc := NewEmailVerificationService(repo, userRepo, &mockEmailSender{}, "http://localhost:5173")

	if err := svc.VerifyToken("used-token"); !errors.Is(err, ErrTokenInvalid) {
		t.Errorf("ожидали ErrTokenInvalid для уже использованного токена, получили: %v", err)
	}
}

func TestVerifyToken_Unknown(t *testing.T) {
	svc := NewEmailVerificationService(&mockEmailVerificationRepo{}, &mockUserRepo{}, &mockEmailSender{}, "http://localhost:5173")

	if err := svc.VerifyToken("nonexistent"); !errors.Is(err, ErrTokenInvalid) {
		t.Errorf("ожидали ErrTokenInvalid для неизвестного токена, получили: %v", err)
	}
}
