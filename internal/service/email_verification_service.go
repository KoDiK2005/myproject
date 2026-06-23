package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"myproject/internal/models"
	"time"
)

var ErrTokenInvalid = errors.New("invalid or expired token")

type EmailVerificationRepository interface {
	Create(token *models.EmailVerificationToken) error
	GetByToken(token string) (*models.EmailVerificationToken, error)
	MarkUsed(token string) error
}

// EmailSender — отправка писем. Реализуется internal/mailer (SMTP или консольный вывод для dev).
type EmailSender interface {
	Send(to, subject, body string) error
}

type EmailVerificationService struct {
	repo        EmailVerificationRepository
	userRepo    UserRepository
	sender      EmailSender
	frontendURL string
}

func NewEmailVerificationService(repo EmailVerificationRepository, userRepo UserRepository, sender EmailSender, frontendURL string) *EmailVerificationService {
	return &EmailVerificationService{repo: repo, userRepo: userRepo, sender: sender, frontendURL: frontendURL}
}

// SendVerification — генерирует токен и отправляет письмо со ссылкой подтверждения.
// Ошибка отправки не должна ронять регистрацию — вызывающий код решает, насколько критично её залогировать.
func (s *EmailVerificationService) SendVerification(userID int, email string) error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	token := &models.EmailVerificationToken{
		UserID:    userID,
		Token:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.Create(token); err != nil {
		return err
	}

	link := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token.Token)
	body := fmt.Sprintf("Подтвердите свой email, перейдя по ссылке (действует 24 часа):\n%s", link)
	return s.sender.Send(email, "Подтверждение email — myproject", body)
}

// VerifyToken — проверяет токен и помечает email подтверждённым
func (s *EmailVerificationService) VerifyToken(tokenStr string) error {
	token, err := s.repo.GetByToken(tokenStr)
	if err != nil || token == nil {
		return ErrTokenInvalid
	}
	if token.Used || token.ExpiresAt.Before(time.Now()) {
		return ErrTokenInvalid
	}
	if err := s.userRepo.MarkEmailVerified(token.UserID); err != nil {
		return err
	}
	return s.repo.MarkUsed(tokenStr)
}
