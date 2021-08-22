package article

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/internal/dynamo"
	"github.com/ferjmc/cms/pkg/rand"
)

type dynamoRepository struct{}

func (d *dynamoRepository) PutArticle(article *entities.Article) error {
	const maxAttempt = 5

	// Try to find a unique article id
	for attempt := 0; ; attempt++ {
		err := putArticleWithRandomId(article)

		if err == nil {
			return nil
		}

		if attempt >= maxAttempt {
			return err
		}

		if !dynamo.IsConditionalCheckFailed(err) {
			return err
		}

		rand.ArticleIdRand.RenewSeed()
	}
}

func putArticleWithRandomId(article *entities.Article) error {
	article.ArticleId = 1 + rand.ArticleIdRand.Get().Int63n(entities.MaxArticleId-1) // range: [1, MaxArticleId)
	article.MakeSlug()

	articleItem, err := dynamodbattribute.MarshalMap(article)
	if err != nil {
		return err
	}

	transactItems := make([]*dynamodb.TransactWriteItem, 0, 1+2*len(article.TagList))

	// Put a new article
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:           aws.String(dynamo.ArticleTableName),
			Item:                articleItem,
			ConditionExpression: aws.String("attribute_not_exists(ArticleId)"),
		},
	})

	for _, tag := range article.TagList {
		articleTag := entities.ArticleTag{
			Tag:       tag,
			ArticleId: article.ArticleId,
			CreatedAt: article.CreatedAt,
		}

		item, err := dynamodbattribute.MarshalMap(articleTag)
		if err != nil {
			return err
		}

		// Link article with tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(dynamo.ArticleTagTableName),
				Item:      item,
			},
		})

		// Update article count for each tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName:        aws.String(dynamo.TagTableName),
				Key:              dynamo.StringKey("Tag", tag),
				UpdateExpression: aws.String("ADD ArticleCount :one SET Dummy=:zero"),
				ExpressionAttributeValues: dynamo.AWSObject{
					":one":  dynamo.IntValue(1),
					":zero": dynamo.IntValue(0),
				},
			},
		})
	}

	_, err = dynamo.DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	return err
}

func (d *dynamoRepository) GetAllArticles(offset, limit int) ([]entities.Article, error) {
	queryArticles := dynamodb.QueryInput{
		TableName:                 aws.String(dynamo.ArticleTableName),
		IndexName:                 aws.String("CreatedAt"),
		KeyConditionExpression:    aws.String("Dummy=:zero"),
		ExpressionAttributeValues: dynamo.IntKey(":zero", 0),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := dynamo.QueryItems(&queryArticles, offset, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]entities.Article, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (d *dynamoRepository) GetArticlesByAuthor(author string, offset, limit int) ([]entities.Article, error) {
	queryArticles := dynamodb.QueryInput{
		TableName:                 aws.String(dynamo.ArticleTableName),
		IndexName:                 aws.String("Author"),
		KeyConditionExpression:    aws.String("Author=:author"),
		ExpressionAttributeValues: dynamo.StringKey(":author", author),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := dynamo.QueryItems(&queryArticles, offset, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]entities.Article, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func GetArticleIdsByTag(tag string, offset, limit int) ([]int64, error) {
	queryArticleIds := dynamodb.QueryInput{
		TableName:                 aws.String(dynamo.ArticleTagTableName),
		IndexName:                 aws.String("CreatedAt"),
		KeyConditionExpression:    aws.String("Tag=:tag"),
		ExpressionAttributeValues: dynamo.StringKey(":tag", tag),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
		ProjectionExpression:      aws.String("ArticleId"),
	}

	items, err := dynamo.QueryItems(&queryArticleIds, offset, limit)
	if err != nil {
		return nil, err
	}

	articleTags := make([]entities.ArticleTag, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articleTags)
	if err != nil {
		return nil, err
	}

	articleIds := make([]int64, 0, len(items))

	for _, articleTag := range articleTags {
		articleIds = append(articleIds, articleTag.ArticleId)
	}

	return articleIds, nil
}

func (d *dynamoRepository) GetArticlesByTag(tag string, offset, limit int) ([]entities.Article, error) {
	articleIds, err := GetArticleIdsByTag(tag, offset, limit)
	if err != nil {
		return nil, err
	}

	return d.GetArticlesByArticleIds(articleIds, limit)
}

func GetFavoriteArticleIdsByUsername(username string, offset, limit int) ([]int64, error) {
	queryArticleIds := dynamodb.QueryInput{
		TableName:                 aws.String(dynamo.FavoriteArticleTableName),
		IndexName:                 aws.String("FavoritedAt"),
		KeyConditionExpression:    aws.String("Username=:username"),
		ExpressionAttributeValues: dynamo.StringKey(":username", username),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
		ProjectionExpression:      aws.String("ArticleId"),
	}

	items, err := dynamo.QueryItems(&queryArticleIds, offset, limit)
	if err != nil {
		return nil, err
	}

	favoriteArticles := make([]entities.FavoriteArticle, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &favoriteArticles)
	if err != nil {
		return nil, err
	}

	articleIds := make([]int64, 0, len(items))

	for _, favoriteArticle := range favoriteArticles {
		articleIds = append(articleIds, favoriteArticle.ArticleId)
	}

	return articleIds, nil
}

func (d *dynamoRepository) GetFavoriteArticlesByUsername(username string, offset, limit int) ([]entities.Article, error) {
	articleIds, err := GetFavoriteArticleIdsByUsername(username, offset, limit)
	if err != nil {
		return nil, err
	}

	return d.GetArticlesByArticleIds(articleIds, limit)
}

func (d *dynamoRepository) GetArticlesByArticleIds(articleIds []int64, limit int) ([]entities.Article, error) {
	if len(articleIds) == 0 {
		return make([]entities.Article, 0), nil
	}

	keys := make([]dynamo.AWSObject, 0, len(articleIds))
	for _, articleId := range articleIds {
		keys = append(keys, dynamo.Int64Key("ArticleId", articleId))
	}

	batchGetArticles := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			dynamo.ArticleTableName: {
				Keys: keys,
			},
		},
	}

	responses, err := dynamo.BatchGetItems(&batchGetArticles, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]entities.Article, len(articleIds))
	articleIdToIndex := dynamo.ReverseIndexInt64(articleIds)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				article := entities.Article{}
				err = dynamodbattribute.UnmarshalMap(item, &article)
				if err != nil {
					return nil, err
				}

				index := articleIdToIndex[article.ArticleId]
				articles[index] = article
			}
		}
	}

	return articles, nil
}
