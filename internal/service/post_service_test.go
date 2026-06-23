package service

import (
	"errors"
	"myproject/internal/models"
	"testing"
)

// мок-репозиторий — подделка которая ничего не знает о БД
type mockPostRepo struct {
	posts   []models.Post
	friends map[[2]int]bool
}

func (m *mockPostRepo) IsFriend(userA, userB int) (bool, error) {
	if m.friends == nil {
		return false, nil
	}
	return m.friends[[2]int{userA, userB}] || m.friends[[2]int{userB, userA}], nil
}

func (m *mockPostRepo) Create(post *models.Post) error {
	post.ID = len(m.posts) + 1
	m.posts = append(m.posts, *post)
	return nil
}

func (m *mockPostRepo) GetByID(id int) (*models.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (m *mockPostRepo) GetAllWithCount(limit, offset int) ([]models.Post, int, error) {
	end := offset + limit
	if end > len(m.posts) {
		end = len(m.posts)
	}
	if offset > len(m.posts) {
		return nil, len(m.posts), nil
	}
	return m.posts[offset:end], len(m.posts), nil
}

func (m *mockPostRepo) GetByUserIDForViewer(ownerID, viewerID, limit, offset int) ([]models.Post, error) {
	var result []models.Post
	for _, p := range m.posts {
		if p.UserID == ownerID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPostRepo) GetFeedWithCount(userID, limit, offset int) ([]models.Post, int, error) {
	return m.GetAllWithCount(limit, offset) // в тестах не имитируем фильтр по дружбе
}

func (m *mockPostRepo) Update(id int, title, body, visibility string) (*models.Post, error) {
	for i, p := range m.posts {
		if p.ID == id {
			m.posts[i].Title = title
			m.posts[i].Body = body
			m.posts[i].Visibility = visibility
			return &m.posts[i], nil
		}
	}
	return nil, nil
}

func (m *mockPostRepo) Delete(id int) error {
	for i, p := range m.posts {
		if p.ID == id {
			m.posts = append(m.posts[:i], m.posts[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockPostRepo) SearchWithCount(query string, limit, offset int) ([]models.Post, int, error) {
	var matched []models.Post
	for _, p := range m.posts {
		if containsStr(p.Title, query) || containsStr(p.Body, query) {
			matched = append(matched, p)
		}
	}
	total := len(matched)
	end := offset + limit
	if end > len(matched) {
		end = len(matched)
	}
	if offset > len(matched) {
		return nil, total, nil
	}
	return matched[offset:end], total, nil
}

func containsStr(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// сам тест
func TestCreatePost(t *testing.T) {
	repo := &mockPostRepo{}
	svc := NewPostService(repo)

	input := models.CreatePostInput{Title: "тест", Body: "тело"}
	post, err := svc.CreatePost(input, 1)

	if err != nil {
		t.Fatalf("не ожидали ошибку, получили: %v", err)
	}
	if post.Title != "тест" {
		t.Errorf("ожидали title='тест', получили '%s'", post.Title)
	}
	if post.UserID != 1 {
		t.Errorf("ожидали user_id=1, получили %d", post.UserID)
	}
}

func TestUpdatePost_Ownership(t *testing.T) {
	repo := &mockPostRepo{}
	svc := NewPostService(repo)

	post, _ := svc.CreatePost(models.CreatePostInput{Title: "old title", Body: "old body"}, 1)

	// чужой юзер — forbidden
	_, err := svc.UpdatePost(post.ID, 2, models.UpdatePostInput{Title: "hack", Body: "hacked"})
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("ожидали ErrForbidden, получили: %v", err)
	}

	// владелец — обновляет нормально
	updated, err := svc.UpdatePost(post.ID, 1, models.UpdatePostInput{Title: "new title", Body: "new body"})
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if updated.Title != "new title" {
		t.Errorf("ожидали 'new title', получили '%s'", updated.Title)
	}
}

func TestDeletePost_Ownership(t *testing.T) {
	repo := &mockPostRepo{}
	svc := NewPostService(repo)

	// создаём пост для юзера 1
	post, _ := svc.CreatePost(models.CreatePostInput{Title: "my post", Body: "body"}, 1)

	// юзер 2 пытается удалить — должен получить forbidden
	err := svc.DeletePost(post.ID, 2)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("ожидали ErrForbidden, получили: %v", err)
	}

	// юзер 1 удаляет свой же пост — всё ок
	err = svc.DeletePost(post.ID, 1)
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
}

func TestGetPostsByUserID(t *testing.T) {
	repo := &mockPostRepo{}
	svc := NewPostService(repo)

	// создай 3 поста — 2 для юзера 1, 1 для юзера 2
	svc.CreatePost(models.CreatePostInput{Title: "post1", Body: "body1"}, 1)
	svc.CreatePost(models.CreatePostInput{Title: "post2", Body: "body2"}, 1)
	svc.CreatePost(models.CreatePostInput{Title: "post3", Body: "body3"}, 2)

	// получи посты юзера 1
	resp, err := svc.GetPostsByUserID(1, 0, 1, 10)

	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	posts := resp.Data.([]models.Post)
	if len(posts) != 2 {
		t.Errorf("ожидали 2 поста, получили %d", len(posts))
	}
}

func TestGetPostByID_VisibilityFriends(t *testing.T) {
	repo := &mockPostRepo{friends: map[[2]int]bool{{1, 2}: true}}
	svc := NewPostService(repo)

	post, _ := svc.CreatePost(models.CreatePostInput{Title: "secret", Body: "body", Visibility: "friends"}, 1)

	// гость не видит приватный пост
	if _, err := svc.GetPostByID(post.ID, 0); !errors.Is(err, ErrNotFound) {
		t.Errorf("гость: ожидали ErrNotFound, получили: %v", err)
	}

	// чужой не-друг (юзер 99) не видит
	if _, err := svc.GetPostByID(post.ID, 99); !errors.Is(err, ErrNotFound) {
		t.Errorf("не-друг: ожидали ErrNotFound, получили: %v", err)
	}

	// друг (юзер 2) видит
	if _, err := svc.GetPostByID(post.ID, 2); err != nil {
		t.Errorf("друг: не ожидали ошибку, получили: %v", err)
	}

	// автор видит свой пост
	if _, err := svc.GetPostByID(post.ID, 1); err != nil {
		t.Errorf("автор: не ожидали ошибку, получили: %v", err)
	}
}
