package service

import "myproject/internal/models"

type CommentRepository interface {
	Create(c *models.Comment) error
	GetByPostID(postID int) ([]models.Comment, error)
	GetByID(id int) (*models.Comment, error)
	Delete(id int) error
}

// PostAccessChecker — проверка видимости поста для пользователя (реализует *PostService)
type PostAccessChecker interface {
	CanViewPost(postID, viewerID int) (bool, error)
}

type CommentService struct {
	repo        CommentRepository
	postChecker PostAccessChecker
}

func NewCommentService(repo CommentRepository, postChecker PostAccessChecker) *CommentService {
	return &CommentService{repo: repo, postChecker: postChecker}
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
	err = s.repo.Create(c)
	return c, err
}

func (s *CommentService) ListComments(postID, viewerID int) ([]models.Comment, error) {
	ok, err := s.postChecker.CanViewPost(postID, viewerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	return s.repo.GetByPostID(postID)
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
