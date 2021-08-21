package user

import (
	"github.com/ferjmc/cms/entities"
)

type mockUserRepository struct{}

func (m *mockUserRepository) PutUser(user entities.User) error {
	return nil
}

func NewMockUserRepository() UserRepository {
	return &mockUserRepository{}
}
