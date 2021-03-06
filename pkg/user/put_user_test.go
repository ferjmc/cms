package user

import (
	"testing"

	"github.com/ferjmc/cms/entities"
)

func TestPutUser(t *testing.T) {
	repo := NewMockUserRepository()
	serv := NewUserService(repo)

	t.Run("It must return an error with a password in blank", func(t *testing.T) {
		user := entities.User{
			Username: "username",
			Email:    "email@fake.com",
		}
		err := serv.PutUser(user, "")
		if err == nil {
			t.Error("get nil error when password blank")
		}
	})

	t.Run("It must accept minimum parameters without errors", func(t *testing.T) {
		user := entities.User{
			Username: "ferjmc",
			Email:    "fernando.castro@telco.com.ar",
		}
		err := serv.PutUser(user, "123456")
		if err != nil {
			t.Errorf("error must be nil, instead: %s", err)
		}
	})
}

func TestGetUserByUsername(t *testing.T) {
	repo := NewMockUserRepository()
	serv := NewUserService(repo)

	t.Run("Given username string retrieve user object with same username atribute", func(t *testing.T) {
		username := "username"
		user, err := serv.GetUserByUsername(username)
		if err != nil {
			t.Error("if repository found the username it can't return an error")
		}
		if username != user.Username {
			t.Errorf("given username: %s, user.username: %s are not equal", username, user.Username)
		}
	})
	t.Run("It must return error if username is empty", func(t *testing.T) {
		_, err := serv.GetUserByUsername("")
		if err == nil {
			t.Error("error is nil while username is blank")
		}
	})
}
