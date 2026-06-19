package handler_test

import (
	"bytes"
	"encoding/json"
	"myproject/internal/handler"
	"myproject/internal/models"
	"myproject/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockPostRepo реализует service.PostRepository без БД
type mockPostRepo struct {
	posts []models.Post
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
	if offset >= len(m.posts) {
		return nil, len(m.posts), nil
	}
	return m.posts[offset:end], len(m.posts), nil
}

func (m *mockPostRepo) GetByUserID(userID, limit, offset int) ([]models.Post, error) {
	var result []models.Post
	for _, p := range m.posts {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPostRepo) GetFeedWithCount(userID, limit, offset int) ([]models.Post, int, error) {
	return m.GetAllWithCount(limit, offset)
}

func (m *mockPostRepo) Update(id int, title, body, visibility string) (*models.Post, error) {
	for i, p := range m.posts {
		if p.ID == id {
			m.posts[i].Title = title
			m.posts[i].Body = body
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
	return nil, 0, nil
}

// setupPostRouter — gin роутер с хэндлером постов
func setupPostRouter(repo service.PostRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewPostService(repo)
	h := handler.NewPostHandler(svc)

	r.GET("/posts", h.ListPosts)
	r.GET("/posts/:id", h.GetPost)
	// user_id задаём через middleware-заглушку
	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})
	authorized.POST("/posts", h.CreatePost)
	authorized.DELETE("/posts/:id", h.DeletePost)

	return r
}

func TestListPosts_EmptyFeed(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("ожидали 200, получили %d", w.Code)
	}
	var resp models.PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("невалидный JSON: %v", err)
	}
	if resp.Total != 0 {
		t.Errorf("ожидали total=0, получили %d", resp.Total)
	}
}

func TestListPosts_WithData(t *testing.T) {
	repo := &mockPostRepo{
		posts: []models.Post{
			{ID: 1, Title: "Первый", Body: "тело", UserID: 1},
			{ID: 2, Title: "Второй", Body: "тело2", UserID: 2},
		},
	}
	r := setupPostRouter(repo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("ожидали 200, получили %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 2 {
		t.Errorf("ожидали total=2, получили %d", resp.Total)
	}
}

func TestListPosts_InvalidPage(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?page=бля", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("ожидали 400 на невалидной странице, получили %d", w.Code)
	}
}

func TestGetPost_NotFound(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("ожидали 404, получили %d", w.Code)
	}
}

func TestGetPost_Found(t *testing.T) {
	repo := &mockPostRepo{
		posts: []models.Post{
			{ID: 1, Title: "Тест", Body: "тело", UserID: 1},
		},
	}
	r := setupPostRouter(repo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("ожидали 200, получили %d", w.Code)
	}
	var post models.Post
	json.Unmarshal(w.Body.Bytes(), &post)
	if post.Title != "Тест" {
		t.Errorf("ожидали title='Тест', получили '%s'", post.Title)
	}
}

func TestCreatePost_Success(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	body, _ := json.Marshal(map[string]string{"title": "Новый пост", "body": "достаточно длинный текст для теста"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("ожидали 201, получили %d: %s", w.Code, w.Body.String())
	}
	var post models.Post
	json.Unmarshal(w.Body.Bytes(), &post)
	if post.UserID != 1 {
		t.Errorf("ожидали user_id=1 (из мидлваре), получили %d", post.UserID)
	}
}

func TestCreatePost_MissingTitle(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	body, _ := json.Marshal(map[string]string{"body": "достаточно длинный текст без заголовка"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("ожидали 400 без title, получили %d", w.Code)
	}
}

func TestDeletePost_Owner_Success(t *testing.T) {
	repo := &mockPostRepo{
		posts: []models.Post{
			{ID: 1, Title: "Мой пост", Body: "тело", UserID: 1}, // owner = user_id 1
		},
	}
	r := setupPostRouter(repo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("ожидали 204, получили %d", w.Code)
	}
}

func TestDeletePost_NotOwner_Forbidden(t *testing.T) {
	repo := &mockPostRepo{
		posts: []models.Post{
			{ID: 1, Title: "Чужой пост", Body: "тело", UserID: 42}, // owner = 42, а мы под user_id=1
		},
	}
	r := setupPostRouter(repo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("ожидали 403, получили %d", w.Code)
	}
}

func TestDeletePost_NotFound(t *testing.T) {
	r := setupPostRouter(&mockPostRepo{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("ожидали 404, получили %d", w.Code)
	}
}
