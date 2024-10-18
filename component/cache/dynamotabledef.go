package cache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
)

func (ac *appCache) GetDynamoDBTableDefinitions() []DynamoDBTableDefinition {
	var enrolledCacheTableDefinition = DynamoDBTableDefinition{
		Name: viper.GetString("DYNAMODB_ENROLLMENT_TABLE_NAME"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("cachekey"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("courseslug"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("cachekey"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("courseslug"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ReadCapacity:     1,
		WriteCapacity:    1,
		TTLAttributeName: "ttl",
	}

	var userCacheTableDefinition = DynamoDBTableDefinition{
		Name: viper.GetString("DYNAMODB_USER_TABLE_NAME"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("cachekey"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("userid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("cachekey"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("userid"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ReadCapacity:     1,
		WriteCapacity:    1,
		TTLAttributeName: "ttl",
	}

	var videoCacheTableDefinition = DynamoDBTableDefinition{
		Name: viper.GetString("DYNAMODB_VIDEO_TABLE_NAME"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("cachekey"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("videoid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("cachekey"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("videoid"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ReadCapacity:     1,
		WriteCapacity:    1,
		TTLAttributeName: "ttl",
	}

	return []DynamoDBTableDefinition{
		enrolledCacheTableDefinition,
		userCacheTableDefinition,
		videoCacheTableDefinition,
	}
}
