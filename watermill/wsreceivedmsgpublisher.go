package watermill

import (
	"context"
	"fmt"
	"salon_be/appconst"
	"salon_be/component/logger"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PublishReceivedWSMsgEvent(
	ctx context.Context,
	localPub *gochannel.GoChannel,
	wsmessage []byte,
) error {
	// Create a Watermill message
	watermillMsg := message.NewMessage(uuid.NewString(), wsmessage)

	// Set the context in the message metadata
	watermillMsg.Metadata.Set("correlation_id", getCorrelationID(ctx))

	userID, ok := ctx.Value("currentUserID").(string)
	if !ok {
		// Handle the case where the value isn't a string
		return fmt.Errorf("userID not found in context")
	}

	watermillMsg.Metadata.Set("currentUserID", userID)
	// Set tracing metadata
	setTracingMetadata(ctx, watermillMsg)

	err := localPub.Publish(appconst.TopicReceivedWSMsg, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(ctx,
			fmt.Sprintf("Error publish %s", appconst.TopicReceivedWSMsg),
			zap.ByteString("msg payload", wsmessage),
			zap.Error(err),
		)
		return err
	}

	return nil
}
