package config

import (
	"fmt"
	"log"
	"os"

	"salon_be/component"
	"salon_be/component/appqueue"
	"salon_be/component/cache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateAWSSession() (*session.Session, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return sess, nil
}

func CreateAppCache(awsSess *session.Session) component.AppCache {
	appcache, err := cache.NewAppCache(awsSess)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	return appcache
}

func CreateAppQueue(awsSession *session.Session) component.AppQueue {
	return appqueue.CreateAppQueue(awsSession)
}

func CreateS3Client(awsSession *session.Session) *s3.S3 {
	return s3.New(awsSession)
}
