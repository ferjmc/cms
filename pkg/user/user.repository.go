package user

import (
	"errors"

	"github.com/ferjmc/cms/entities"
)

const (
	InstanceDynamodb int = iota
)

type UserRepository interface {
	PutUser(user entities.User) error
}

func NewUserRepository(instance int) (UserRepository, error) {
	switch instance {
	case InstanceDynamodb:
		return &dynamoRepository{}, nil
	default:
		return nil, errors.New("repository instance not found")
	}
}
