package appqueue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type appQ struct {
	sqsClient *sqs.SQS
}

func CreateAppQueue(awsSess *session.Session) *appQ {
	return &appQ{
		sqsClient: sqs.New(awsSess),
	}
}
