package service

import "myproject/internal/models"

type BlockRepository interface {
	Block(blockerID, blockedID int) error
	Unblock(blockerID, blockedID int) error
	IsBlocked(userA, userB int) (bool, error)
	GetBlockedUsers(blockerID int) ([]models.User, error)
}

// FriendRemover — узкий интерфейс, реализуется FriendshipRepository.
// Блокировка должна разрывать существующую дружбу/заявки между юзерами.
type FriendRemover interface {
	Remove(userID, friendID int) error
}

type BlockService struct {
	repo    BlockRepository
	friends FriendRemover
}

func NewBlockService(repo BlockRepository, friends FriendRemover) *BlockService {
	return &BlockService{repo: repo, friends: friends}
}

func (s *BlockService) Block(blockerID, blockedID int) error {
	if blockerID == blockedID {
		return ErrForbidden
	}
	if err := s.repo.Block(blockerID, blockedID); err != nil {
		return err
	}
	return s.friends.Remove(blockerID, blockedID)
}

func (s *BlockService) Unblock(blockerID, blockedID int) error {
	return s.repo.Unblock(blockerID, blockedID)
}

func (s *BlockService) IsBlocked(userA, userB int) (bool, error) {
	return s.repo.IsBlocked(userA, userB)
}

func (s *BlockService) GetBlockedUsers(blockerID int) ([]models.User, error) {
	return s.repo.GetBlockedUsers(blockerID)
}
