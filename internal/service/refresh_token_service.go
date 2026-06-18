package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"myproject/internal/models"
	"time"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	GetByToken(token string) (*models.RefreshToken, error)
	Revoke(token string) error
}

type RefreshTokenService struct {
	repo RefreshTokenRepository
}

func NewRefreshTokenService(repo RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{repo: repo}
}

func (s *RefreshTokenService) GenerateToken(userID int) (*models.RefreshToken, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	token := &models.RefreshToken{
		UserID:    userID,
		Token:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // живёт 7 дней
	}

	err := s.repo.Create(token)
	return token, err
}

func (s *RefreshTokenService) ValidateToken(token string) (*models.RefreshToken, error) {
	rt, err := s.repo.GetByToken(token)
	if err != nil {
		return nil, ErrNotFound
	}
	if rt.Revoked {
		return nil, errors.New("token revoked")
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	return rt, nil
}

func (s *RefreshTokenService) RevokeToken(token string) error {
	return s.repo.Revoke(token)
}
