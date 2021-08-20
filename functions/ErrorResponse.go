package functions

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ferjmc/cms/entities"
)

type InputErrorResponse struct {
	Errors entities.InputError `json:"errors"`
}

func NewErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	inputError, ok := err.(entities.InputError)
	if !ok {
		// Internal server error
		return events.APIGatewayProxyResponse{}, err
	}

	body := InputErrorResponse{
		Errors: inputError,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 422,
		Body:       string(jsonBody),
		Headers:    CORSHeaders(),
	}
	return response, nil
}

func NewUnauthorizedResponse() (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: 401,
		Headers:    CORSHeaders(),
	}
	return response, nil
}
