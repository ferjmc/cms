package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ferjmc/cms/functions"
	"github.com/ferjmc/cms/pkg/follow"
	"github.com/ferjmc/cms/pkg/user"
)

type Response struct {
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	Username  string `json:"username"`
	Image     string `json:"image"`
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userService := user.New()
	user, _, err := userService.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return functions.NewUnauthorizedResponse()
	}

	publisher, err := userService.GetUserByUsername(input.PathParameters["username"])
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	followService := follow.New()
	err = followService.Unfollow(user.Username, publisher.Username)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	response := Response{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: false,
		},
	}

	return functions.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
