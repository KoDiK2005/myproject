package service

type LikeRepository interface {
	Like(userID, postID int) error
	Unlike(userID, postID int) error
	Count(postID int) (int, error)
	IsLiked(userID, postID int) (bool, error)
}

type LikeService struct {
	repo LikeRepository
}

func NewLikeService(repo LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) Like(userID, postID int) error {
	return s.repo.Like(userID, postID)
}

func (s *LikeService) Unlike(userID, postID int) error {
	return s.repo.Unlike(userID, postID)
}

func (s *LikeService) Count(postID int) (int, error) {
	return s.repo.Count(postID)
}

func (s *LikeService) IsLiked(userID, postID int) (bool, error) {
	if userID == 0 {
		return false, nil
	}
	return s.repo.IsLiked(userID, postID)
}
