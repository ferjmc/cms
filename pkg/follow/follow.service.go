package follow

import "github.com/ferjmc/cms/entities"

type FollowService interface {
	IsFollowing(follower *entities.User, publishers []string) ([]bool, error)
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
	return s.repository.IsFollowing(follower, publishers)
}
