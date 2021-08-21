package user

import (
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/pkg/auth"
)

type UserService interface {
	PutUser(user entities.User, password string) error
	GetUserByUsername(username string) (*entities.User, error)
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

func (s *userService) GetUserByUsername(username string) (*entities.User, error) {
	if len(username) <= 0 {
		return nil, entities.NewInputError("username", "username can't be blank")
	}
	return nil, nil
}
