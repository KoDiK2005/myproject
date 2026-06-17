package service

import "myproject/internal/models"

type CommentRepository interface {
	Create(c *models.Comment) error
	GetByPostID(postID int) ([]models.Comment, error)
	GetByID(id int) (*models.Comment, error)
	Delete(id int) error
}

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(postID, userID int, input models.CreateCommentInput) (*models.Comment, error) {
	c := &models.Comment{
		PostID: postID,
		UserID: userID,
		Body:   input.Body,
	}
	err := s.repo.Create(c)
	return c, err
}

func (s *CommentService) ListComments(postID int) ([]models.Comment, error) {
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
