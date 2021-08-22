package follow

import "github.com/ferjmc/cms/entities"

type FollowService interface {
	// IsFollowing given a user and a list of publishers, retrieves a list
	// with a flag in true for each publisher if the user is following his posts
	IsFollowing(follower *entities.User, publishers []string) ([]bool, error)
	Follow(follower, publisher string) error
	Unfollow(follower, publisher string) error
}

func NewFollowService(r FollowRepository) FollowService {
	return &followService{
		repository: r,
	}
}

func New(opts ...func(FollowService) FollowService) FollowService {
	var serv FollowService
	for _, opt := range opts {
		serv = opt(serv)
	}
	// whitout opts retrieves service with dynamo by default
	if len(opts) <= 0 {
		return WithDynamoDB(serv)
	}
	return serv
}

func WithDynamoDB(serv FollowService) FollowService {
	repo, err := NewFollowRepository(InstanceDynamodb)
	if err != nil {
		return serv
	}
	return NewFollowService(repo)
}

type followService struct {
	repository FollowRepository
}

func (s *followService) IsFollowing(follower *entities.User, publishers []string) ([]bool, error) {
	if follower == nil || len(publishers) == 0 {
		return make([]bool, len(publishers)), nil
	}
	return s.repository.IsFollowing(follower, publishers)
}

func (s *followService) Follow(follower, publisher string) error {
	follow := entities.Follow{
		Follower:  follower,
		Publisher: publisher,
	}
	return s.repository.Follow(follow)
}

func (s *followService) Unfollow(follower, publisher string) error {
	follow := entities.Follow{
		Follower:  follower,
		Publisher: publisher,
	}
	return s.repository.Unfollow(follow)
}
