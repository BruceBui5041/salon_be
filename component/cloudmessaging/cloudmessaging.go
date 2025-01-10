package cloudmessaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"salon_be/component/logger"
	models "salon_be/model"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type fcmClient struct {
	client *messaging.Client
}

func NewFCMClient(credentialPath string) (*fcmClient, error) {
	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	messagingClient, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &fcmClient{
		client: messagingClient,
	}, nil
}

func (f *fcmClient) SendNotification(ctx context.Context, notification *models.Notification) error {
	if notification == nil {
		return errors.New("notification cannot be nil")
	}

	// Check if notification should be sent now or scheduled
	if notification.Scheduled != nil && notification.Scheduled.After(time.Now()) {
		return nil // Will be sent by scheduler when time comes
	}

	// Convert notification metadata to map
	metadata := make(map[string]string)
	if notification.Metadata != nil {
		for k, v := range notification.Metadata {
			if str, ok := v.(string); ok {
				metadata[k] = str
			} else {
				// Convert complex values to JSON string
				if jsonStr, err := json.Marshal(v); err == nil {
					metadata[k] = string(jsonStr)
				}
			}
		}
	}

	jsonMetadata, err := json.Marshal(notification.Metadata)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"Error marshalling notification metadata",
			zap.Error(err),
			zap.Any("notification", notification),
		)
		return err
	}

	// Create FCM message
	message := &messaging.Message{
		Notification: &messaging.Notification{
			// Title: notification.Title,
			Body: string(jsonMetadata),
		},
		Data: metadata,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Sound:             "default",
				NotificationCount: pointer(1),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound:    "default",
					Badge:    pointer(1),
					Category: notification.Type,
				},
			},
		},
	}

	// Send to each recipient in notification details
	for _, detail := range notification.Details {
		if detail.User.UserDevice.FCMToken != "" {
			message.Token = detail.User.UserDevice.FCMToken

			// Send message
			result, err := f.client.Send(ctx, message)
			if err != nil {
				detail.Status = models.NotificationStateError
				detail.Error = err.Error()
			} else {
				detail.Status = models.NotificationStateSent
				detail.MessageID = result
				detail.SentAt = time.Now().UTC()
			}
		}
	}

	return nil
}

func (f *fcmClient) SendBatchNotifications(ctx context.Context, notifications []*models.Notification) error {
	for _, notification := range notifications {
		if err := f.SendNotification(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

// Helper function to get pointer to integer
func pointer(i int) *int {
	return &i
}
