package service

import (
	"errors"
	"myproject/internal/models"
	"testing"
	"time"
)

type mockCommentRepo struct {
	comments []models.Comment
}

func (m *mockCommentRepo) Create(c *models.Comment) error {
	c.ID = len(m.comments) + 1
	c.CreatedAt = time.Now()
	m.comments = append(m.comments, *c)
	return nil
}

func (m *mockCommentRepo) GetByPostID(postID int) ([]models.Comment, error) {
	var result []models.Comment
	for _, c := range m.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCommentRepo) GetByID(id int) (*models.Comment, error) {
	for _, c := range m.comments {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockCommentRepo) Delete(id int) error {
	for i, c := range m.comments {
		if c.ID == id {
			m.comments = append(m.comments[:i], m.comments[i+1:]...)
			return nil
		}
	}
	return nil
}

// allowAllPostChecker — мок PostAccessChecker, который всегда разрешает доступ
type allowAllPostChecker struct{}

func (allowAllPostChecker) CanViewPost(postID, viewerID int) (bool, error) {
	return true, nil
}

func TestAddComment(t *testing.T) {
	repo := &mockCommentRepo{}
	svc := NewCommentService(repo, allowAllPostChecker{})

	c, err := svc.AddComment(1, 2, models.CreateCommentInput{Body: "привет"})
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if c.PostID != 1 || c.UserID != 2 || c.Body != "привет" {
		t.Errorf("неверные данные комментария: %+v", c)
	}
	if c.ID == 0 {
		t.Error("ID не проставился")
	}
}

func TestListComments(t *testing.T) {
	repo := &mockCommentRepo{
		comments: []models.Comment{
			{ID: 1, PostID: 10, UserID: 1, Body: "первый"},
			{ID: 2, PostID: 10, UserID: 2, Body: "второй"},
			{ID: 3, PostID: 99, UserID: 1, Body: "не наш пост"},
		},
	}
	svc := NewCommentService(repo, allowAllPostChecker{})

	comments, err := svc.ListComments(10, 0)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("ожидали 2 комментария для post 10, получили %d", len(comments))
	}
}

func TestDeleteComment_Owner(t *testing.T) {
	repo := &mockCommentRepo{
		comments: []models.Comment{
			{ID: 1, PostID: 1, UserID: 5, Body: "мой коммент"},
		},
	}
	svc := NewCommentService(repo, allowAllPostChecker{})

	err := svc.DeleteComment(1, 5) // user 5 удаляет свой коммент
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if len(repo.comments) != 0 {
		t.Error("коммент должен быть удалён")
	}
}

func TestDeleteComment_NotOwner(t *testing.T) {
	repo := &mockCommentRepo{
		comments: []models.Comment{
			{ID: 1, PostID: 1, UserID: 5, Body: "чужой коммент"},
		},
	}
	svc := NewCommentService(repo, allowAllPostChecker{})

	err := svc.DeleteComment(1, 99) // user 99 лезет в чужой коммент
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("ожидали ErrForbidden, получили: %v", err)
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	svc := NewCommentService(&mockCommentRepo{}, allowAllPostChecker{})

	err := svc.DeleteComment(999, 1)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("ожидали ErrNotFound, получили: %v", err)
	}
}
