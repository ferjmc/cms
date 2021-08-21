package user

import (
	"github.com/ferjmc/cms/entities"
)

type mockUserRepository struct{}

func (m *mockUserRepository) PutUser(user entities.User) error {
	return nil
}

func (m *mockUserRepository) UserByUsername(username string) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserRepository) UsernameByEmail(email string) (string, error) {
	return "", nil
}

func (m *mockUserRepository) UpdateUser(oldUser, newUser entities.User) error {
	return nil
}

func NewMockUserRepository() UserRepository {
	return &mockUserRepository{}
}
