package user

import (
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/pkg/auth"
)

type UserService interface {
	// PutUser creates a new user from basic struct and a password string
	PutUser(user entities.User, password string) error
	// GetUserByUsername retrieves a user object from a username string
	GetUserByUsername(username string) (*entities.User, error)
	GetUsernameByEmail(email string) (string, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetCurrentUser(authorization string) (*entities.User, string, error)
	UpdateUser(authorization string, newUser entities.User) (*entities.User, string, error)
	GetUserListByUsername(usernames []string) ([]entities.User, error)
}

func NewUserService(r UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

func New(opts ...func(UserService) UserService) UserService {
	var serv UserService
	for _, opt := range opts {
		serv = opt(serv)
	}
	// whitout opts retrieves service with dynamo by default
	if len(opts) <= 0 {
		return WithDynamoDB(serv)
	}
	return serv
}

func WithDynamoDB(serv UserService) UserService {
	repo, err := NewUserRepository(InstanceDynamodb)
	if err != nil {
		return serv
	}
	return NewUserService(repo)
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

func (s *userService) GetUserListByUsername(usernames []string) ([]entities.User, error) {
	if len(usernames) == 0 {
		return make([]entities.User, 0), nil
	}

	return s.repository.GetUserListByUsername(usernames)
}
