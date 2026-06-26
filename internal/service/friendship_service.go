package service

import "myproject/internal/models"

type FriendshipRepository interface {
	SendRequest(requesterID, addresseeID int) (bool, error)
	Accept(requesterID, addresseeID int) error
	Reject(requesterID, addresseeID int) error
	Remove(userID, friendID int) error
	GetFriends(userID int) ([]models.User, error)
	GetIncomingRequests(userID int) ([]models.User, error)
	GetOutgoingRequests(userID int) ([]models.User, error)
	GetStatus(userID, otherID int) (*models.FriendshipStatus, error)
}

// BlockChecker — узкий интерфейс, реализуется BlockService.
type BlockChecker interface {
	IsBlocked(userA, userB int) (bool, error)
}

type FriendshipService struct {
	repo     FriendshipRepository
	blocks   BlockChecker
	notifier Notifier
}

func NewFriendshipService(repo FriendshipRepository, blocks BlockChecker, notifier Notifier) *FriendshipService {
	return &FriendshipService{repo: repo, blocks: blocks, notifier: notifier}
}

func (s *FriendshipService) SendRequest(requesterID, addresseeID int) error {
	if requesterID == addresseeID {
		return ErrForbidden // себе заявку не отправить
	}
	blocked, err := s.blocks.IsBlocked(requesterID, addresseeID)
	if err != nil {
		return err
	}
	if blocked {
		return ErrForbidden
	}
	autoAccepted, err := s.repo.SendRequest(requesterID, addresseeID)
	if err != nil {
		return err
	}
	if autoAccepted {
		// у addresseeID уже была своя исходящая заявка к requesterID — она была принята
		_ = s.notifier.Notify(addresseeID, requesterID, "friend_accept", nil)
	} else {
		_ = s.notifier.Notify(addresseeID, requesterID, "friend_request", nil)
	}
	return nil
}

func (s *FriendshipService) Accept(currentUserID, requesterID int) error {
	if err := s.repo.Accept(requesterID, currentUserID); err != nil {
		return err
	}
	_ = s.notifier.Notify(requesterID, currentUserID, "friend_accept", nil)
	return nil
}

func (s *FriendshipService) Reject(currentUserID, requesterID int) error {
	return s.repo.Reject(requesterID, currentUserID)
}

func (s *FriendshipService) Remove(userID, friendID int) error {
	return s.repo.Remove(userID, friendID)
}

func (s *FriendshipService) GetFriends(userID int) ([]models.User, error) {
	return s.repo.GetFriends(userID)
}

func (s *FriendshipService) GetIncomingRequests(userID int) ([]models.User, error) {
	return s.repo.GetIncomingRequests(userID)
}

func (s *FriendshipService) GetOutgoingRequests(userID int) ([]models.User, error) {
	return s.repo.GetOutgoingRequests(userID)
}

func (s *FriendshipService) GetStatus(userID, otherID int) (*models.FriendshipStatus, error) {
	return s.repo.GetStatus(userID, otherID)
}
