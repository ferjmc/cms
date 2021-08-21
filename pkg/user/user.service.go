package user

import (
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/pkg/auth"
)

type UserService interface {
	PutUser(user entities.User, password string) error
	GetUserByUsername(username string) (*entities.User, error)
	GetUsernameByEmail(email string) (string, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetCurrentUser(authorization string) (*entities.User, string, error)
	UpdateUser(authorization string, newUser entities.User) (*entities.User, string, error)
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
	return s.repository.UserByUsername(username)
}

func (s *userService) GetUsernameByEmail(email string) (string, error) {
	return s.repository.UsernameByEmail(email)
}

func (s *userService) GetUserByEmail(email string) (*entities.User, error) {
	username, err := s.repository.UsernameByEmail(email)
	if err != nil {
		return nil, err
	}
	return s.repository.UserByUsername(username)
}

func (s *userService) GetCurrentUser(authorization string) (*entities.User, string, error) {
	authService := auth.New()
	username, token, err := authService.VerifyAuthorization(authorization)
	if err != nil {
		return nil, "", err
	}
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) UpdateUser(authorization string, newUser entities.User) (*entities.User, string, error) {
	err := newUser.Validate()
	if err != nil {
		return nil, "", err
	}

	oldUser, token, err := s.GetCurrentUser(authorization)
	if err != nil {
		return nil, "", err
	}

	err = s.repository.UpdateUser(*oldUser, newUser)
	if err != nil {
		return nil, "", err
	}

	return &newUser, token, nil
}
