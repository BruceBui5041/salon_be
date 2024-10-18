package watermill

import (
	"encoding/json"
	"video_server/component"
	"video_server/component/logger"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func HandleVideoProcessed(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "HandleVideoProcessed")
	defer span.End()

	var processStateInfo *messagemodel.VideoProcessStateInfo
	err := json.Unmarshal(msg.Payload, &processStateInfo)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"Cannot unmarshal message payload",
			zap.Any("payload", msg.Payload),
			zap.Error(err),
		)
		return
	}

	msg.Ack()
}
