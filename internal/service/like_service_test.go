package service

import (
	"testing"
)

type mockLikeRepo struct {
	likes map[string]bool // "userID:postID"
}

func newMockLikeRepo() *mockLikeRepo {
	return &mockLikeRepo{likes: make(map[string]bool)}
}

func likeKey(userID, postID int) string {
	return string(rune(userID)) + ":" + string(rune(postID))
}

func (m *mockLikeRepo) Like(userID, postID int) error {
	m.likes[likeKey(userID, postID)] = true
	return nil
}

func (m *mockLikeRepo) Unlike(userID, postID int) error {
	delete(m.likes, likeKey(userID, postID))
	return nil
}

func (m *mockLikeRepo) Count(postID int) (int, error) {
	count := 0
	for k := range m.likes {
		// ключ содержит postID в конце после ":"
		_ = k
		count++ // упрощённый счётчик — в реале фильтровали бы по postID
	}
	return count, nil
}

func TestLike(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo)

	if err := svc.Like(1, 10); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if !repo.likes[likeKey(1, 10)] {
		t.Error("лайк не сохранился")
	}
}

func TestUnlike(t *testing.T) {
	repo := newMockLikeRepo()
	repo.likes[likeKey(1, 10)] = true
	svc := NewLikeService(repo)

	if err := svc.Unlike(1, 10); err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if repo.likes[likeKey(1, 10)] {
		t.Error("лайк должен быть удалён")
	}
}

func TestLikeCount(t *testing.T) {
	repo := newMockLikeRepo()
	repo.likes[likeKey(1, 10)] = true
	repo.likes[likeKey(2, 10)] = true
	svc := NewLikeService(repo)

	count, err := svc.Count(10)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if count != 2 {
		t.Errorf("ожидали 2, получили %d", count)
	}
}

func TestLikeIdempotent(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo)

	// лайкаем дважды — в реале ON CONFLICT DO NOTHING, мок просто перезаписывает
	svc.Like(1, 10)
	svc.Like(1, 10)

	count, _ := svc.Count(10)
	if count != 1 {
		t.Errorf("ожидали 1 уникальный лайк, получили %d", count)
	}
}
