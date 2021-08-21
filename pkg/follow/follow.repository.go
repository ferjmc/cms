package follow

import (
	"errors"

	"github.com/ferjmc/cms/entities"
)

const (
	InstanceDynamodb int = iota
)

type FollowRepository interface {
	IsFollowing(follower *entities.User, publishers []string) ([]bool, error)
	Follow(follow entities.Follow) error
}

func NewFollowRepository(instance int) (FollowRepository, error) {
	switch instance {
	case InstanceDynamodb:
		return &dynamoRepository{}, nil
	default:
		return nil, errors.New("repository instance not found")
	}
}
