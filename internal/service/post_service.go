package service

import (
	"myproject/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id int) (*models.Post, error)
	GetAllWithCount(limit, offset int) ([]models.Post, int, error)
	GetFeedWithCount(userID, limit, offset int) ([]models.Post, int, error)
	GetByUserIDForViewer(ownerID, viewerID, limit, offset int) ([]models.Post, error)
	Update(id int, title, body, visibility string) (*models.Post, error)
	Delete(id int) error
	SearchWithCount(query string, viewerID, limit, offset int) ([]models.Post, int, error)
	IsFriend(userA, userB int) (bool, error)
	IsBlocked(userA, userB int) (bool, error)
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

// GetPostByID — viewerID=0 для гостя. Приватные посты (friends) видят только
// автор и его друзья; остальным возвращаем ErrNotFound, чтобы не подтверждать
// существование поста.
func (s *PostService) GetPostByID(id, viewerID int) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil || post == nil {
		return nil, ErrNotFound
	}
	ok, err := s.canView(post, viewerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	return post, nil
}

// CanViewPost — то же самое, но без выдачи поста; используется другими
// сервисами (например комментариями) для проверки доступа.
func (s *PostService) CanViewPost(postID, viewerID int) (bool, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil || post == nil {
		return false, nil
	}
	return s.canView(post, viewerID)
}

func (s *PostService) canView(post *models.Post, viewerID int) (bool, error) {
	if viewerID != 0 && viewerID != post.UserID {
		blocked, err := s.repo.IsBlocked(post.UserID, viewerID)
		if err != nil {
			return false, err
		}
		if blocked {
			return false, nil
		}
	}
	if post.Visibility != "friends" || viewerID == post.UserID {
		return true, nil
	}
	if viewerID == 0 {
		return false, nil
	}
	return s.repo.IsFriend(post.UserID, viewerID)
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

// GetPostsByUserID — посты конкретного юзера с учётом того кто смотрит
func (s *PostService) GetPostsByUserID(ownerID, viewerID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, err := s.repo.GetByUserIDForViewer(ownerID, viewerID, limit, offset)
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

func (s *PostService) SearchPosts(query string, viewerID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.SearchWithCount(query, viewerID, limit, offset)
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
