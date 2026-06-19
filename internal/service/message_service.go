package service

import (
	"errors"
	"myproject/internal/models"
)

var ErrNotFriends = errors.New("you can only message friends")

type MessageRepository interface {
	Save(msg *models.Message) error
	GetHistory(userA, userB, limit, offset int) ([]models.Message, error)
	GetConversations(userID int) ([]models.ConversationPreview, error)
	MarkRead(receiverID, senderID int) error
	IsFriend(userA, userB int) (bool, error)
}

type MessageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// Send — отправить сообщение. Проверяем дружбу перед сохранением.
func (s *MessageService) Send(senderID, receiverID int, content string) (*models.Message, error) {
	ok, err := s.repo.IsFriend(senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFriends
	}
	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	if err := s.repo.Save(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageService) GetHistory(userA, userB, page, limit int) ([]models.Message, error) {
	offset := (page - 1) * limit
	return s.repo.GetHistory(userA, userB, limit, offset)
}

func (s *MessageService) GetConversations(userID int) ([]models.ConversationPreview, error) {
	return s.repo.GetConversations(userID)
}

func (s *MessageService) MarkRead(receiverID, senderID int) error {
	return s.repo.MarkRead(receiverID, senderID)
}
