package service

import "myproject/internal/models"

type NotificationRepository interface {
	Create(n *models.Notification) error
	GetByUserID(userID, limit, offset int) ([]models.Notification, int, error)
	MarkRead(userID, notificationID int) error
	MarkAllRead(userID int) error
	UnreadCount(userID int) (int, error)
}

// NotificationPusher — доставка уведомления онлайн-юзеру по WebSocket (реализует *ws.Hub)
type NotificationPusher interface {
	DeliverNotification(n *models.Notification)
}

// Notifier — узкий интерфейс, реализуется *NotificationService.
// Используется другими сервисами (friendship, like, comment) чтобы создавать уведомления
// без прямой зависимости от репозитория уведомлений.
type Notifier interface {
	Notify(userID, actorID int, notifType string, postID *int) error
}

type NotificationService struct {
	repo   NotificationRepository
	pusher NotificationPusher
}

func NewNotificationService(repo NotificationRepository, pusher NotificationPusher) *NotificationService {
	return &NotificationService{repo: repo, pusher: pusher}
}

// Notify — создать уведомление и доставить онлайн-юзеру. Не уведомляем самого себя
// (например лайк собственного поста не создаёт уведомление).
func (s *NotificationService) Notify(userID, actorID int, notifType string, postID *int) error {
	if userID == actorID {
		return nil
	}
	n := &models.Notification{
		UserID:  userID,
		ActorID: actorID,
		Type:    notifType,
		PostID:  postID,
	}
	if err := s.repo.Create(n); err != nil {
		return err
	}
	s.pusher.DeliverNotification(n)
	return nil
}

func (s *NotificationService) GetByUserID(userID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	list, total, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *NotificationService) MarkRead(userID, notificationID int) error {
	return s.repo.MarkRead(userID, notificationID)
}

func (s *NotificationService) MarkAllRead(userID int) error {
	return s.repo.MarkAllRead(userID)
}

func (s *NotificationService) UnreadCount(userID int) (int, error) {
	return s.repo.UnreadCount(userID)
}
