package service

import (
	"errors"
	"myproject/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id int) (*models.Post, error)
	GetAll(limit, offset int) ([]models.Post, error)
	GetByUserID(userID, limit, offset int) ([]models.Post, error)
	Delete(id int) error
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(input models.CreatePostInput, userID int) (*models.Post, error) {
	post := &models.Post{
		Title:  input.Title,
		Body:   input.Body,
		UserID: userID,
	}
	err := s.repo.Create(post)
	return post, err
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) ListPosts(page, limit int) ([]models.Post, error) {
	offset := (page - 1) * limit
	return s.repo.GetAll(limit, offset)
}
func (s *PostService) GetPostsByUserID(userID, page, limit int) ([]models.Post, error) {
	offset := (page - 1) * limit
	return s.repo.GetByUserID(userID, limit, offset)
}
func (s *PostService) DeletePost(id, userID int) error {
	post, err := s.repo.GetByID(id)
	if err != nil || post == nil {
		return errors.New("post not found")
	}
	if post.UserID != userID {
		// чужой пост — нечего тут делать
		return errors.New("forbidden")
	}
	return s.repo.Delete(id)
}
