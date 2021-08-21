package user

import (
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/pkg/auth"
)

type UserService interface {
	PutUser(user entities.User, password string) error
}

func NewUserService(r UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

type userService struct {
	repository UserRepository
}

func (s *userService) PutUser(user entities.User, password string) error {
	err := entities.ValidatePassword(password)
	if err != nil {
		return err
	}

	passHash, err := auth.New().Scrypt(password)
	if err != nil {
		return err
	}

	user.PasswordHash = passHash

	err = user.Validate()
	if err != nil {
		return err
	}

	return s.repository.PutUser(user)
}
