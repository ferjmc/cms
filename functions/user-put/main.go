package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/functions"
	"github.com/ferjmc/cms/pkg/auth"
	"github.com/ferjmc/cms/pkg/user"
)

type Request struct {
	User UserRequest `json:"user"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

type Response struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo, err := user.NewUserRepository(user.InstanceDynamodb)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	var request Request
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return functions.NewErrorResponse(err)
	}
	auth := auth.New()
	passwordHash, err := auth.Scrypt(request.User.Password)
	if err != nil {
		return functions.NewErrorResponse(err)
	}
	newUser := entities.User{
		Username:     request.User.Username,
		Email:        request.User.Email,
		PasswordHash: passwordHash,
		Image:        request.User.Image,
		Bio:          request.User.Bio,
	}

	userService := user.NewUserService(repo)
	user, token, err := userService.UpdateUser(input.Headers["Authorization"], newUser)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	response := Response{
		User: UserResponse{
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
			Bio:      user.Bio,
			Token:    token,
		},
	}

	return functions.NewSuccessResponse(200, response)
}
