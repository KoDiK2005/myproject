package service

import "myproject/internal/models"

type CommentRepository interface {
	Create(c *models.Comment) error
	GetByPostID(postID, limit, offset int) ([]models.Comment, error)
	CountByPostID(postID int) (int, error)
	GetByID(id int) (*models.Comment, error)
	Delete(id int) error
}

// PostAccessChecker — проверка видимости поста и владельца (реализует *PostService)
type PostAccessChecker interface {
	CanViewPost(postID, viewerID int) (bool, error)
	GetOwnerID(postID int) (int, error)
}

type CommentService struct {
	repo        CommentRepository
	postChecker PostAccessChecker
	notifier    Notifier
}

func NewCommentService(repo CommentRepository, postChecker PostAccessChecker, notifier Notifier) *CommentService {
	return &CommentService{repo: repo, postChecker: postChecker, notifier: notifier}
}

func (s *CommentService) AddComment(postID, userID int, input models.CreateCommentInput) (*models.Comment, error) {
	ok, err := s.postChecker.CanViewPost(postID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	c := &models.Comment{
		PostID: postID,
		UserID: userID,
		Body:   input.Body,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	if ownerID, err := s.postChecker.GetOwnerID(postID); err == nil {
		_ = s.notifier.Notify(ownerID, userID, "comment", &postID)
	}
	return c, nil
}

func (s *CommentService) ListComments(postID, viewerID, page, limit int) (*models.PaginatedResponse, error) {
	ok, err := s.postChecker.CanViewPost(postID, viewerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	offset := (page - 1) * limit
	comments, err := s.repo.GetByPostID(postID, limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountByPostID(postID)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: comments, Total: total, Page: page, Limit: limit}, nil
}

func (s *CommentService) DeleteComment(id, userID int) error {
	c, err := s.repo.GetByID(id)
	if err != nil || c == nil {
		return ErrNotFound
	}
	if c.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}
