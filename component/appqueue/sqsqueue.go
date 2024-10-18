package appqueue

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"video_server/component"
	"video_server/component/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type SQSDefinition struct {
	QueueName              string
	Delay                  time.Duration
	MaximumMessageSize     int64 // in bytes
	MessageRetentionPeriod time.Duration
	ReceiveMessageWaitTime time.Duration
	VisibilityTimeout      time.Duration
}

func (appQ *appQ) generateSQSQueueURL(queueName string) string {
	region := viper.GetString("AWS_REGION")
	accountID := viper.GetString("AWS_ACCOUNT_ID")

	return fmt.Sprintf(
		"https://sqs.%s.amazonaws.com/%s/%s",
		region,
		accountID,
		queueName,
	)
}

func (appQ *appQ) StartSQSMessageListener(
	ctx context.Context,
	appContext component.AppContext,
	consumeTopics []string,
	processMessageHandler func(ctx context.Context, appContext component.AppContext, msg *sqs.Message) error,
) {
	// Start polling for messages in separate goroutines for each topic
	for _, topic := range consumeTopics {
		go func(t string, appContext component.AppContext) {
			logger.AppLogger.Info(ctx, "Starting to poll messages", zap.String("topic", t))
			appQ.PollSQSMessages(ctx, appContext, t, processMessageHandler)
		}(topic, appContext)
	}
}

func (appQ *appQ) SendSQSMessage(ctx context.Context, topic, groupId, messageBody string) error {
	queueURL := appQ.generateSQSQueueURL(topic)
	logger.AppLogger.Info(ctx, "SendSQSMessage", zap.String("topic", queueURL), zap.String("body", messageBody))
	_, err := appQ.sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:       aws.String(queueURL),
		MessageBody:    aws.String(messageBody),
		MessageGroupId: aws.String(groupId),
	})
	return err
}

func (appQ *appQ) ReceiveSQSMessage(ctx context.Context, topic string) ([]*sqs.Message, error) {
	queueURL := appQ.generateSQSQueueURL(topic)
	result, err := appQ.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
	})
	if err != nil {
		return nil, err
	}
	return result.Messages, nil
}

func (appQ *appQ) PollSQSMessages(
	ctx context.Context,
	appContext component.AppContext,
	topic string,
	handler func(context.Context, component.AppContext, *sqs.Message) error,
) {
	queueURL := appQ.generateSQSQueueURL(topic)
	for {
		messages, err := appQ.ReceiveSQSMessage(ctx, topic)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to receive messages",
				zap.String("topic", topic),
				zap.Error(err))
			continue
		}

		for _, msg := range messages {
			err := handler(ctx, appContext, msg)
			if err != nil {
				logger.AppLogger.Error(ctx, "Failed to process message",
					zap.String("topic", topic),
					zap.Error(err))
				continue
			}

			// Delete the message from the queue after successful processing
			_, err = appQ.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				logger.AppLogger.Error(ctx, "Failed to delete message",
					zap.String("topic", topic),
					zap.Error(err))
			}
		}

		// If no messages were received, add a small delay to avoid tight looping
		if len(messages) == 0 {
			time.Sleep(time.Second)
		}
	}
}

func (a *appQ) createSQSQueue(ctx context.Context, def SQSDefinition) (*string, error) {
	// Check if queue already exists
	existingURL, err := a.getQueueURL(def.QueueName)
	if err != nil {
		return nil, fmt.Errorf("error checking queue existence: %w", err)
	}
	if existingURL != nil {
		return existingURL, nil // Queue already exists, return its URL
	}

	// Validate queue name
	if err := validateSQSQueueName(def.QueueName); err != nil {
		return nil, fmt.Errorf("invalid queue name: %w", err)
	}

	// Prepare queue attributes
	attributes := map[string]*string{
		"DelaySeconds":                  aws.String(strconv.FormatInt(int64(def.Delay.Seconds()), 10)),
		"MaximumMessageSize":            aws.String(strconv.FormatInt(def.MaximumMessageSize, 10)),
		"MessageRetentionPeriod":        aws.String(strconv.FormatInt(int64(def.MessageRetentionPeriod.Seconds()), 10)),
		"ReceiveMessageWaitTimeSeconds": aws.String(strconv.FormatInt(int64(def.ReceiveMessageWaitTime.Seconds()), 10)),
		"VisibilityTimeout":             aws.String(strconv.FormatInt(int64(def.VisibilityTimeout.Seconds()), 10)),
	}

	// Check if it's a FIFO queue
	isFIFO := strings.HasSuffix(def.QueueName, ".fifo")
	if isFIFO {
		attributes["FifoQueue"] = aws.String("true")
		attributes["ContentBasedDeduplication"] = aws.String("true")
	}

	// Create the queue
	input := &sqs.CreateQueueInput{
		QueueName:  aws.String(def.QueueName),
		Attributes: attributes,
	}

	result, err := a.sqsClient.CreateQueue(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	return result.QueueUrl, nil
}

// CreateSQSQueues creates multiple SQS queues concurrently, checking for existence first
func (a *appQ) CreateSQSQueues(ctx context.Context) map[string]*string {
	defs := a.getSQSDefinitions()

	results := make(map[string]*string)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, def := range defs {
		wg.Add(1)
		go func(def SQSDefinition) {
			defer wg.Done()
			queueURL, err := a.createSQSQueue(ctx, def) // Pass context here
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				logger.AppLogger.Fatal(ctx, "failed to create new sqs", zap.Error(err))
			} else {
				results[def.QueueName] = queueURL
			}
		}(def)
	}

	wg.Wait()
	return results
}

func (a *appQ) getSQSDefinitions() []SQSDefinition {
	return []SQSDefinition{
		{
			QueueName:              viper.GetString("REQ_PROCESS_VIDEO_TOPIC"),
			Delay:                  0 * time.Second,
			MaximumMessageSize:     262144,             // 256 KiB
			MessageRetentionPeriod: 4 * 24 * time.Hour, // 4 days
			ReceiveMessageWaitTime: 0 * time.Second,
			VisibilityTimeout:      time.Hour,
		},
		{
			QueueName:              viper.GetString("UPDATE_VIDEO_PROCESS_PROGRESS_TOPIC"),
			Delay:                  0 * time.Second,
			MaximumMessageSize:     131072,             // 128 KiB
			MessageRetentionPeriod: 3 * 24 * time.Hour, // 3 days
			ReceiveMessageWaitTime: 0 * time.Second,
			VisibilityTimeout:      time.Hour,
		},
	}
}

func (a *appQ) getQueueURL(queueName string) (*string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := a.sqsClient.GetQueueUrl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			return nil, nil // Queue doesn't exist
		}
		return nil, err // Some other error occurred
	}

	return result.QueueUrl, nil
}

func validateSQSQueueName(name string) error {
	if !strings.HasSuffix(name, ".fifo") {
		if len(name) < 1 || len(name) > 80 {
			return fmt.Errorf("standard queue name must be between 1 and 80 characters long")
		}
		validName := regexp.MustCompile(`^[0-9A-Za-z-_]+$`)
		if !validName.MatchString(name) {
			return fmt.Errorf("standard queue name can only include alphanumeric characters, hyphens, or underscores")
		}
	} else {
		// FIFO queue name validation
		nameWithoutSuffix := strings.TrimSuffix(name, ".fifo")
		if len(nameWithoutSuffix) < 1 || len(nameWithoutSuffix) > 75 {
			return fmt.Errorf("FIFO queue name must be between 1 and 75 characters long (excluding the .fifo suffix)")
		}
		validName := regexp.MustCompile(`^[0-9A-Za-z-_]+\.fifo$`)
		if !validName.MatchString(name) {
			return fmt.Errorf("FIFO queue name can only include alphanumeric characters, hyphens, or underscores, and must end with .fifo")
		}
	}
	return nil
}
