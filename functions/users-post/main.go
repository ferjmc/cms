package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	request := Request{}
	err := json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	repo, err := user.NewUserRepository(user.InstanceDynamodb)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	service := user.NewUserService(repo)

	newUser := entities.User{
		Username:     request.User.Username,
		Email:        request.User.Email,
		PasswordHash: []byte(request.User.Password),
	}

	err = service.PutUser(newUser)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	token, err := auth.New().GenerateToken(newUser.Username)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	response := Response{
		User: UserResponse{
			Username: newUser.Username,
			Email:    newUser.Email,
			Image:    newUser.Image,
			Bio:      newUser.Bio,
			Token:    token,
		},
	}

	return functions.NewSuccessResponse(201, response)
}

func main() {
	lambda.Start(Handle)
}
