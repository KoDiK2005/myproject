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
	Update(id int, title, body string) (*models.Post, error)
	Count() (int, error)
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

func (s *PostService) ListPosts(page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.Count()
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
	// считаем только посты этого юзера
	total := len(posts)
	return &models.PaginatedResponse{Data: posts, Total: total, Page: page, Limit: limit}, nil
}
func (s *PostService) UpdatePost(id, userID int, input models.UpdatePostInput) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}
	if post.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return s.repo.Update(id, input.Title, input.Body)
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
