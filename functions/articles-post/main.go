package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/functions"
	"github.com/ferjmc/cms/pkg/article"
	"github.com/ferjmc/cms/pkg/user"
)

type Request struct {
	Article ArticleRequest `json:"article"`
}

type ArticleRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type Response struct {
	Article ArticleResponse `json:"article"`
}

type ArticleResponse struct {
	Slug           string         `json:"slug"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Body           string         `json:"body"`
	TagList        []string       `json:"tagList"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
	Favorited      bool           `json:"favorited"`
	FavoritesCount int64          `json:"favoritesCount"`
	Author         AuthorResponse `json:"author"`
}

type AuthorResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := user.New().GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return functions.NewUnauthorizedResponse()
	}

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(entities.TimestampFormat)

	newArticle := entities.Article{
		Title:       request.Article.Title,
		Description: request.Article.Description,
		Body:        request.Article.Body,
		TagList:     request.Article.TagList, // TODO .distinct()
		CreatedAt:   nowUnixNano,
		UpdatedAt:   nowUnixNano,
		Author:      user.Username,
	}

	err = article.New().PutArticle(&newArticle)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	response := Response{
		Article: ArticleResponse{
			Slug:           newArticle.Slug,
			Title:          newArticle.Title,
			Description:    newArticle.Description,
			Body:           newArticle.Body,
			TagList:        newArticle.TagList,
			CreatedAt:      nowStr,
			UpdatedAt:      nowStr,
			Favorited:      false,
			FavoritesCount: 0,
			Author: AuthorResponse{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: false,
			},
		},
	}

	return functions.NewSuccessResponse(201, response)
}

func main() {
	lambda.Start(Handle)
}
