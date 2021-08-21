package dynamo

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var once sync.Once
var svc *dynamodb.DynamoDB

func initializeSingletons() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
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
	return fmt.Sprintf("cms-%s-%s", Stage, suffix)
}
