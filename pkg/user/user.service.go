package user

import "github.com/ferjmc/cms/entities"

type UserService interface {
	PutUser(user entities.User) error
}

func NewUserService(r UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

type userService struct {
	repository UserRepository
}

func (s *userService) PutUser(user entities.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	return s.repository.PutUser(user)
}
