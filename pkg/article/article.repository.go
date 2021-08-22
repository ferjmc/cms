package article

import (
	"errors"

	"github.com/ferjmc/cms/entities"
)

const (
	InstanceDynamodb int = iota
)

type ArticleRepository interface {
	PutArticle(article *entities.Article) error
	GetAllArticles(offset, limit int) ([]entities.Article, error)
	GetArticlesByAuthor(author string, offset, limit int) ([]entities.Article, error)
	GetArticlesByTag(tag string, offset, limit int) ([]entities.Article, error)
	GetFavoriteArticlesByUsername(username string, offset, limit int) ([]entities.Article, error)
}

func NewArticleRepository(instance int) (ArticleRepository, error) {
	switch instance {
	case InstanceDynamodb:
		return &dynamoRepository{}, nil
	default:
		return nil, errors.New("repository instance not found")
	}
}
