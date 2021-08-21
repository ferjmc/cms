package follow

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/internal/dynamo"
)

type dynamoRepository struct{}

func (d *dynamoRepository) IsFollowing(follower *entities.User, publishers []string) ([]bool, error) {
	publisherSet := make(map[string]bool)
	for _, publisher := range publishers {
		publisherSet[publisher] = true
	}

	keys := make([]dynamo.AWSObject, 0, len(publisherSet))
	for publisher := range publisherSet {
		keys = append(keys, dynamo.AWSObject{
			"Follower":  dynamo.StringValue(follower.Username),
			"Publisher": dynamo.StringValue(publisher),
		})
	}

	batchGetFollows := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			dynamo.FollowTableName: {
				Keys:                 keys,
				ProjectionExpression: aws.String("Publisher"),
			},
		},
	}

	responses, err := dynamo.BatchGetItems(&batchGetFollows, len(publisherSet))
	if err != nil {
		return nil, err
	}

	followingUser := make(map[string]bool)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				var follow entities.Follow
				err = dynamodbattribute.UnmarshalMap(item, &follow)
				if err != nil {
					return nil, err
				}
				followingUser[follow.Publisher] = true
			}
		}
	}

	following := make([]bool, 0, len(publishers))
	for _, username := range publishers {
		following = append(following, followingUser[username])
	}

	return following, nil
}

func (d *dynamoRepository) Follow(follow entities.Follow) error {
	item, err := dynamodbattribute.MarshalMap(follow)
	if err != nil {
		return err
	}

	putFollow := dynamodb.PutItemInput{
		TableName: aws.String(dynamo.FollowTableName),
		Item:      item,
	}

	_, err = dynamo.DynamoDB().PutItem(&putFollow)

	return err
}
