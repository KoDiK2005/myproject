package service

type LikeRepository interface {
	Like(userID, postID int) error
	Unlike(userID, postID int) error
	Count(postID int) (int, error)
	IsLiked(userID, postID int) (bool, error)
}

type LikeService struct {
	repo        LikeRepository
	postChecker PostAccessChecker
	notifier    Notifier
}

func NewLikeService(repo LikeRepository, postChecker PostAccessChecker, notifier Notifier) *LikeService {
	return &LikeService{repo: repo, postChecker: postChecker, notifier: notifier}
}

func (s *LikeService) Like(userID, postID int) error {
	ok, err := s.postChecker.CanViewPost(postID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	if err := s.repo.Like(userID, postID); err != nil {
		return err
	}
	if ownerID, err := s.postChecker.GetOwnerID(postID); err == nil {
		_ = s.notifier.Notify(ownerID, userID, "like", &postID)
	}
	return nil
}

func (s *LikeService) Unlike(userID, postID int) error {
	ok, err := s.postChecker.CanViewPost(postID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return s.repo.Unlike(userID, postID)
}

func (s *LikeService) Count(postID, viewerID int) (int, error) {
	ok, err := s.postChecker.CanViewPost(postID, viewerID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrNotFound
	}
	return s.repo.Count(postID)
}

func (s *LikeService) IsLiked(userID, postID int) (bool, error) {
	if userID == 0 {
		return false, nil
	}
	return s.repo.IsLiked(userID, postID)
}
