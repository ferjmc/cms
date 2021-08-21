package main

import (
	"bytes"
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

	serv := user.NewUserService(repo)

	user, err := serv.GetUserByEmail(request.User.Email)
	if err != nil {
		return functions.NewErrorResponse(err)
	}
	auth := auth.New()
	passwordHash, err := auth.Scrypt(request.User.Password)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	if !bytes.Equal(passwordHash, user.PasswordHash) {
		return functions.NewErrorResponse(entities.NewInputError("Password", "password incorrect!"))
	}

	token, err := auth.GenerateToken(user.Username)
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

func main() {
	lambda.Start(Handle)
}
