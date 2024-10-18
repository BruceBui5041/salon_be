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

func PublishVideoUploadedEvent(
	ctx context.Context,
	appCtx component.AppContext,
	videoInfo *messagemodel.RequestProcessVideoInfo,
) error {
	// Marshal videoInfo into JSON
	payload, err := json.Marshal(videoInfo)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"Error marshaling videoInfo to JSON",
			zap.Any("videoInfo", videoInfo),
			zap.Error(err),
		)
		return err
	}

	// Create a Watermill message
	watermillMsg := message.NewMessage(uuid.NewString(), payload)
	// Set tracing metadata
	setTracingMetadata(ctx, watermillMsg)

	err = appCtx.GetLocalPubSub().GetUnblockPubSub().Publish(appconst.TopicNewVideoUploaded, watermillMsg)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			fmt.Sprintf("Error publish %s", appconst.TopicNewVideoUploaded),
			zap.Any("msg payload", payload),
			zap.Error(err),
		)
		return err
	}

	return nil
}
