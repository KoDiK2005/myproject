package service

import (
	"errors"
	"myproject/internal/models"
	"testing"
)

type mockFriendshipRepo struct {
	// "requester:addressee" -> status
	rows map[[2]int]string
}

func newMockFriendshipRepo() *mockFriendshipRepo {
	return &mockFriendshipRepo{rows: make(map[[2]int]string)}
}

func (m *mockFriendshipRepo) SendRequest(requesterID, addresseeID int) (bool, error) {
	if status, ok := m.rows[[2]int{addresseeID, requesterID}]; ok && status == "pending" {
		m.rows[[2]int{addresseeID, requesterID}] = "accepted"
		return true, nil
	}
	if _, ok := m.rows[[2]int{requesterID, addresseeID}]; !ok {
		m.rows[[2]int{requesterID, addresseeID}] = "pending"
	}
	return false, nil
}

func (m *mockFriendshipRepo) Accept(requesterID, addresseeID int) error {
	m.rows[[2]int{requesterID, addresseeID}] = "accepted"
	return nil
}

func (m *mockFriendshipRepo) Reject(requesterID, addresseeID int) error {
	delete(m.rows, [2]int{requesterID, addresseeID})
	return nil
}

func (m *mockFriendshipRepo) Remove(userID, friendID int) error {
	delete(m.rows, [2]int{userID, friendID})
	delete(m.rows, [2]int{friendID, userID})
	return nil
}

func (m *mockFriendshipRepo) GetFriends(userID int) ([]models.User, error) { return nil, nil }
func (m *mockFriendshipRepo) GetIncomingRequests(userID int) ([]models.User, error) {
	return nil, nil
}
func (m *mockFriendshipRepo) GetOutgoingRequests(userID int) ([]models.User, error) {
	return nil, nil
}
func (m *mockFriendshipRepo) GetStatus(userID, otherID int) (*models.FriendshipStatus, error) {
	return nil, nil
}

type allowAllBlockChecker struct{}

func (allowAllBlockChecker) IsBlocked(userA, userB int) (bool, error) { return false, nil }

func TestSendRequest_SelfForbidden(t *testing.T) {
	svc := NewFriendshipService(newMockFriendshipRepo(), allowAllBlockChecker{}, noopNotifier{})

	if err := svc.SendRequest(1, 1); !errors.Is(err, ErrForbidden) {
		t.Errorf("ожидали ErrForbidden, получили: %v", err)
	}
}

func TestSendRequest_NotifiesAddressee(t *testing.T) {
	notifier := &recordingNotifier{}
	svc := NewFriendshipService(newMockFriendshipRepo(), allowAllBlockChecker{}, notifier)

	if err := svc.SendRequest(1, 2); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if !notifier.called || notifier.userID != 2 || notifier.actorID != 1 || notifier.typ != "friend_request" {
		t.Errorf("неверное уведомление: called=%v userID=%d actorID=%d type=%s", notifier.called, notifier.userID, notifier.actorID, notifier.typ)
	}
}

func TestSendRequest_MutualAutoAccepts(t *testing.T) {
	repo := newMockFriendshipRepo()
	notifier := &recordingNotifier{}
	svc := NewFriendshipService(repo, allowAllBlockChecker{}, notifier)

	// 1 отправляет заявку 2
	if err := svc.SendRequest(1, 2); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	// 2 отправляет встречную заявку 1 — должно авто-принять, а не создать вторую pending
	if err := svc.SendRequest(2, 1); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if repo.rows[[2]int{1, 2}] != "accepted" {
		t.Errorf("ожидали accepted, получили %q", repo.rows[[2]int{1, 2}])
	}
	if !notifier.called || notifier.userID != 1 || notifier.actorID != 2 || notifier.typ != "friend_accept" {
		t.Errorf("ожидали friend_accept для userID=1 actorID=2, получили userID=%d actorID=%d type=%s", notifier.userID, notifier.actorID, notifier.typ)
	}
}

func TestAccept_NotifiesRequester(t *testing.T) {
	repo := newMockFriendshipRepo()
	repo.rows[[2]int{1, 2}] = "pending"
	notifier := &recordingNotifier{}
	svc := NewFriendshipService(repo, allowAllBlockChecker{}, notifier)

	if err := svc.Accept(2, 1); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if !notifier.called || notifier.userID != 1 || notifier.actorID != 2 || notifier.typ != "friend_accept" {
		t.Errorf("неверное уведомление: called=%v userID=%d actorID=%d type=%s", notifier.called, notifier.userID, notifier.actorID, notifier.typ)
	}
}

func TestSendRequest_BlockedForbidden(t *testing.T) {
	svc := NewFriendshipService(newMockFriendshipRepo(), denyAllBlockChecker{}, noopNotifier{})

	if err := svc.SendRequest(1, 2); !errors.Is(err, ErrForbidden) {
		t.Errorf("ожидали ErrForbidden, получили: %v", err)
	}
}

type denyAllBlockChecker struct{}

func (denyAllBlockChecker) IsBlocked(userA, userB int) (bool, error) { return true, nil }
