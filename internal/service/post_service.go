package service

import (
	"myproject/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id int) (*models.Post, error)
	GetAllWithCount(limit, offset int) ([]models.Post, int, error)
	GetFeedWithCount(userID, limit, offset int) ([]models.Post, int, error)
	GetByUserID(userID, limit, offset int) ([]models.Post, error)
	Update(id int, title, body, visibility string) (*models.Post, error)
	Delete(id int) error
	SearchWithCount(query string, limit, offset int) ([]models.Post, int, error)
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(input models.CreatePostInput, userID int) (*models.Post, error) {
	v := input.Visibility
	if v != "friends" {
		v = "public" // дефолт — публичный
	}
	post := &models.Post{
		Title:      input.Title,
		Body:       input.Body,
		UserID:     userID,
		Visibility: v,
	}
	err := s.repo.Create(post)
	return post, err
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.GetByID(id)
}

// ListPosts — публичная лента (без авторизации)
func (s *PostService) ListPosts(page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.GetAllWithCount(limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: posts, Total: total, Page: page, Limit: limit}, nil
}

// GetFeed — персональная лента для авторизованного юзера
func (s *PostService) GetFeed(userID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.GetFeedWithCount(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: posts, Total: total, Page: page, Limit: limit}, nil
}

func (s *PostService) GetPostsByUserID(userID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: posts, Total: len(posts), Page: page, Limit: limit}, nil
}

func (s *PostService) UpdatePost(id, userID int, input models.UpdatePostInput) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil || post == nil {
		return nil, ErrNotFound
	}
	if post.UserID != userID {
		return nil, ErrForbidden
	}
	v := input.Visibility
	if v != "friends" {
		v = "public"
	}
	return s.repo.Update(id, input.Title, input.Body, v)
}

func (s *PostService) SearchPosts(query string, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.SearchWithCount(query, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{Data: posts, Total: total, Page: page, Limit: limit}, nil
}

func (s *PostService) DeletePost(id, userID int) error {
	post, err := s.repo.GetByID(id)
	if err != nil || post == nil {
		return ErrNotFound
	}
	if post.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}
