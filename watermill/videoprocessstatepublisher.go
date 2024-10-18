package watermill

import (
	"context"
	"encoding/json"
	"fmt"
	"video_server/appconst"
	"video_server/component"
	"video_server/component/logger"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PublishUpdateProcessVideoStateEvent(
	ctx context.Context,
	appCtx component.AppContext,
	processStateInfo *messagemodel.VideoProcessStateInfo,
) error {
	payload, err := json.Marshal(processStateInfo)
	if err != nil {
		logger.AppLogger.Error(ctx,
			"Error marshaling processStateInfo to JSON",
			zap.Any("processStateInfo", processStateInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(uuid.NewString(), payload)
	// Set tracing metadata
	setTracingMetadata(ctx, watermillMsg)

	// Set the context in the message metadata
	watermillMsg.Metadata.Set("correlation_id", getCorrelationID(ctx))

	err = appCtx.GetLocalPubSub().GetUnblockPubSub().Publish(appconst.TopicUpdateVideoProcessingState, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(ctx,
			fmt.Sprintf("Error publish %s", appconst.TopicUpdateVideoProcessingState),
			zap.Any("msg payload", payload),
			zap.Error(err),
		)
		return err
	}

	return nil
}
