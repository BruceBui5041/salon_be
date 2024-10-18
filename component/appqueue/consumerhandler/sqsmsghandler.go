package consumerhandler

import (
	"context"
	"fmt"
	"video_server/component"
	"video_server/component/logger"

	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

func QueueMsgHander(ctx context.Context, appContext component.AppContext, msg *sqs.Message) error {
	messageGroupID, ok := msg.Attributes[sqs.MessageSystemAttributeNameMessageGroupId]
	logger.AppLogger.Info(ctx, "new msg received", zap.Any("messageGroupID", messageGroupID), zap.Any("msg", msg))

	if !ok {
		fmt.Println("Message group ID not found")
	} else {
		fmt.Printf("Message Group ID: %s\n", *messageGroupID)
	}

	return nil
}
