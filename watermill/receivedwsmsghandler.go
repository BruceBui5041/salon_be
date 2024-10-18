package watermill

import (
	"encoding/json"
	"video_server/component"
	"video_server/component/logger"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func ReceivedWSMsgHandler(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "ReceivedWSMsgHandler")
	defer span.End()

	logger.AppLogger.Info(ctx, "ReceivedWSMsgHandler", zap.Any("msg payload", msg))

	var updateUserCacheInfo *messagemodel.EnrollmentChangeInfo
	err := json.Unmarshal(msg.Payload, &updateUserCacheInfo)
	if err != nil {
		msg.Ack()
		logger.AppLogger.Error(ctx, "Cannot unmarshal message payload", zap.Any("payload", msg.Payload), zap.Error(err))
		return
	}

	msg.Ack()

}
