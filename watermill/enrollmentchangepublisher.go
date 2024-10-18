package watermill

import (
	"context"
	"encoding/json"
	"fmt"
	"salon_be/appconst"
	"salon_be/component/logger"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PublishEnrollmentChange(
	ctx context.Context,
	localPub *gochannel.GoChannel,
	updateCacheInfo *messagemodel.EnrollmentChangeInfo,
) error {
	payload, err := json.Marshal(updateCacheInfo)
	if err != nil {
		logger.AppLogger.Error(ctx,
			"Error marshaling updateCacheInfo to JSON",
			zap.Any("updateCacheInfo", updateCacheInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(uuid.NewString(), payload)

	// Set the context in the message metadata
	watermillMsg.Metadata.Set("correlation_id", getCorrelationID(ctx))
	// Set tracing metadata
	setTracingMetadata(ctx, watermillMsg)

	err = localPub.Publish(appconst.TopicEnrollmentChange, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(ctx,
			fmt.Sprintf("Error publish %s", appconst.TopicEnrollmentChange),
			zap.Any("msg payload", payload),
			zap.Error(err),
		)
		return err
	}

	return nil
}
