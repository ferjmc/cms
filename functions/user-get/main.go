package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ferjmc/cms/functions"
	"github.com/ferjmc/cms/pkg/user"
)

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
	userService := user.NewUserService(repo)
	user, token, err := userService.GetCurrentUser(input.Headers["Authorization"])
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
