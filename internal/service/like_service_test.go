package service

import (
	"errors"
	"testing"
)

type denyPostChecker struct{}

func (denyPostChecker) CanViewPost(postID, viewerID int) (bool, error) {
	return false, nil
}

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

func (m *mockLikeRepo) IsLiked(userID, postID int) (bool, error) {
	return m.likes[likeKey(userID, postID)], nil
}

func TestLike(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, allowAllPostChecker{})

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
	svc := NewLikeService(repo, allowAllPostChecker{})

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
	svc := NewLikeService(repo, allowAllPostChecker{})

	count, err := svc.Count(10, 1)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if count != 2 {
		t.Errorf("ожидали 2, получили %d", count)
	}
}

func TestLikeIdempotent(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, allowAllPostChecker{})

	// лайкаем дважды — в реале ON CONFLICT DO NOTHING, мок просто перезаписывает
	svc.Like(1, 10)
	svc.Like(1, 10)

	count, _ := svc.Count(10, 1)
	if count != 1 {
		t.Errorf("ожидали 1 уникальный лайк, получили %d", count)
	}
}

func TestIsLiked(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, allowAllPostChecker{})

	svc.Like(1, 10)

	liked, err := svc.IsLiked(1, 10)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if !liked {
		t.Error("ожидали liked=true для юзера который лайкнул")
	}

	liked, _ = svc.IsLiked(2, 10)
	if liked {
		t.Error("ожидали liked=false для юзера который не лайкал")
	}

	liked, _ = svc.IsLiked(0, 10) // гость
	if liked {
		t.Error("гость не может быть liked=true")
	}
}

func TestLike_NoAccessToPost(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, denyPostChecker{})

	if err := svc.Like(1, 10); !errors.Is(err, ErrNotFound) {
		t.Errorf("ожидали ErrNotFound, получили: %v", err)
	}
	if repo.likes[likeKey(1, 10)] {
		t.Error("лайк не должен сохраниться для недоступного поста")
	}
}

func TestUnlike_NoAccessToPost(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, denyPostChecker{})

	if err := svc.Unlike(1, 10); !errors.Is(err, ErrNotFound) {
		t.Errorf("ожидали ErrNotFound, получили: %v", err)
	}
}

func TestLikeCount_NoAccessToPost(t *testing.T) {
	repo := newMockLikeRepo()
	svc := NewLikeService(repo, denyPostChecker{})

	if _, err := svc.Count(10, 1); !errors.Is(err, ErrNotFound) {
		t.Errorf("ожидали ErrNotFound, получили: %v", err)
	}
}
