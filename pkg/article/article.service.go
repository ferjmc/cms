package article

import (
	"errors"
	"fmt"

	"github.com/ferjmc/cms/entities"
)

type ArticleService interface {
	PutArticle(article *entities.Article) error
}

func NewArticleService(r ArticleRepository) ArticleService {
	return &articleService{
		repository: r,
	}
}

func New(opts ...func(ArticleService) ArticleService) ArticleService {
	var serv ArticleService
	for _, opt := range opts {
		serv = opt(serv)
	}
	// whitout opts retrieves service with dynamo by default
	if len(opts) <= 0 {
		return WithDynamoDB(serv)
	}
	return serv
}

func WithDynamoDB(serv ArticleService) ArticleService {
	repo, err := NewArticleRepository(InstanceDynamodb)
	if err != nil {
		return serv
	}
	return NewArticleService(repo)
}

type articleService struct {
	repository ArticleRepository
}

func (s *articleService) PutArticle(article *entities.Article) error {
	err := article.Validate()
	if err != nil {
		return err
	}

	return s.repository.PutArticle(article)
}

func (s *articleService) GetArticles(offset, limit int, author, tag, favorited string) ([]entities.Article, error) {
	if offset < 0 {
		return nil, entities.NewInputError("offset", "must be non-negative")
	}

	if limit <= 0 {
		return nil, entities.NewInputError("limit", "must be positive")
	}

	const maxDepth = 1000
	if offset+limit > maxDepth {
		return nil, entities.NewInputError("offset + limit", fmt.Sprintf("must be smaller or equal to %d", maxDepth))
	}

	numFilters := getNumFilters(author, tag, favorited)
	if numFilters > 1 {
		return nil, entities.NewInputError("author, tag, favorited", "only one of these can be specified")
	}

	if numFilters == 0 {
		return s.repository.GetAllArticles(offset, limit)
	}

	if author != "" {
		return s.repository.GetArticlesByAuthor(author, offset, limit)
	}

	if tag != "" {
		return s.repository.GetArticlesByTag(tag, offset, limit)
	}

	if favorited != "" {
		return s.repository.GetFavoriteArticlesByUsername(favorited, offset, limit)
	}

	return nil, errors.New("unreachable code")
}

func getNumFilters(author, tag, favorited string) int {
	numFilters := 0
	if author != "" {
		numFilters++
	}
	if tag != "" {
		numFilters++
	}
	if favorited != "" {
		numFilters++
	}
	return numFilters
}
