package main

import (
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/functions"
	"github.com/ferjmc/cms/pkg/article"
	"github.com/ferjmc/cms/pkg/user"
)

type Response struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
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
	userService := user.New()
	user, _, err := userService.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return functions.NewUnauthorizedResponse()
	}

	offset, err := strconv.Atoi(input.QueryStringParameters["offset"])
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(input.QueryStringParameters["limit"])
	if err != nil {
		limit = 20
	}

	author := input.QueryStringParameters["author"]
	tag := input.QueryStringParameters["tag"]
	favorited := input.QueryStringParameters["favorited"]

	articleService := article.New()
	articles, err := articleService.GetArticles(offset, limit, author, tag, favorited)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	isFavorited, authors, following, err := articleService.GetArticleRelatedProperties(user, articles, true)
	if err != nil {
		return functions.NewErrorResponse(err)
	}

	articleResponses := make([]ArticleResponse, 0, len(articles))

	for i, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(entities.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(entities.TimestampFormat),
			Favorited:      isFavorited[i],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: following[i],
			},
		})
	}

	response := Response{
		Articles:      articleResponses,
		ArticlesCount: len(articleResponses),
	}

	return functions.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
