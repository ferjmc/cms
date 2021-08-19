package dynamo

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var once sync.Once
var svc *dynamodb.DynamoDB

func initializeSingletons() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
}

func DynamoDB() *dynamodb.DynamoDB {
	once.Do(initializeSingletons)
	return svc
}

var Stage = os.Getenv("STAGE")

var UserTableName = makeTableName("user")
var EmailUserTableName = makeTableName("email-user")
var FollowTableName = makeTableName("follow")
var ArticleTableName = makeTableName("article")
var ArticleTagTableName = makeTableName("article-tag")
var TagTableName = makeTableName("tag")
var FavoriteArticleTableName = makeTableName("favorite-article")
var CommentTableName = makeTableName("comment")

func makeTableName(suffix string) string {
	return fmt.Sprintf("realworld-%s-%s", Stage, suffix)
}
