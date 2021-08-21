package user

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ferjmc/cms/entities"
	"github.com/ferjmc/cms/internal/dynamo"
)

type dynamoRepository struct {
}

func (d *dynamoRepository) PutUser(user entities.User) error {
	userItem, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	emailUser := entities.EmailUser{
		Email:    user.Email,
		Username: user.Username,
	}

	emailUserItem, err := dynamodbattribute.MarshalMap(emailUser)
	if err != nil {
		return err
	}

	// Put a new user, make sure username and email are unique
	transaction := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(dynamo.UserTableName),
					Item:                userItem,
					ConditionExpression: aws.String("attribute_not_exists(Username)"),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(dynamo.EmailUserTableName),
					Item:                emailUserItem,
					ConditionExpression: aws.String("attribute_not_exists(Email)"),
				},
			},
		},
	}

	_, err = dynamo.DynamoDB().TransactWriteItems(&transaction)
	if err != nil {
		log.Printf("ERROR: during transaction %s", err)
		// TODO: distinguish:
		// NewInputError("username", "has already been taken")
		// NewInputError("email", "has already been taken")
		return err
	}

	return nil
}

func (d *dynamoRepository) UserByUsername(username string) (*entities.User, error) {
	var user entities.User
	found, err := dynamo.GetItemByKey(dynamo.UserTableName, dynamo.StringKey("Username", username), &user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, entities.NewInputError("username", "not found")
	}

	return &user, nil
}

func (d *dynamoRepository) UsernameByEmail(email string) (string, error) {
	var emailUser entities.EmailUser
	found, err := dynamo.GetItemByKey(dynamo.EmailUserTableName, dynamo.StringKey("Email", email), &emailUser)

	if err != nil {
		return "", err
	}

	if !found {
		return "", entities.NewInputError("email", "not found")
	}

	return emailUser.Username, nil
}

func (d *dynamoRepository) UpdateUser(oldUser, newUser entities.User) error {
	transactItems := make([]*dynamodb.TransactWriteItem, 0, 3)

	if oldUser.Email != newUser.Email {
		newEmailUser := entities.EmailUser{
			Email:    newUser.Email,
			Username: newUser.Username,
		}

		newEmailUserItem, err := dynamodbattribute.MarshalMap(newEmailUser)
		if err != nil {
			return err
		}

		// Link user with the new email
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           aws.String(dynamo.EmailUserTableName),
				Item:                newEmailUserItem,
				ConditionExpression: aws.String("attribute_not_exists(Email)"),
			},
		})

		// Unlink user from the old email
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				TableName:           aws.String(dynamo.EmailUserTableName),
				Key:                 dynamo.StringKey("Email", oldUser.Email),
				ConditionExpression: aws.String("attribute_exists(Email)"),
			},
		})
	}

	newUserItem, err := dynamodbattribute.MarshalMap(newUser)
	if err != nil {
		return err
	}

	// Update user info
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:                 aws.String(dynamo.UserTableName),
			Item:                      newUserItem,
			ConditionExpression:       aws.String("Email = :email"),
			ExpressionAttributeValues: dynamo.StringKey(":email", oldUser.Email),
		},
	})

	_, err = dynamo.DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return err
	}

	return nil
}
