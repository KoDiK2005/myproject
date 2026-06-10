package service

import (
	"myproject/internal/models"
	"myproject/internal/repository"
)

type PostService struct {
	repo *repository.PostRepo
}

func NewPostService(repo *repository.PostRepo) *PostService {
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

func (s *PostService) DeletePost(id int) error {
	return s.repo.Delete(id)
}
