package watermill

import (
	"encoding/json"
	"video_server/component"
	"video_server/component/logger"
	"video_server/model/user/userutils"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func UserUpdatedHandler(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "UserUpdatedHandler")
	defer span.End()

	var updatedUserInfo *messagemodel.UserUpdatedMessage
	err := json.Unmarshal(msg.Payload, &updatedUserInfo)
	if err != nil {
		logger.AppLogger.Error(ctx, "Cannot unmarshal message payload", zap.Any("payload", msg.Payload), zap.Error(err))
		msg.Ack()
		return
	}

	if err := userutils.UpdateUserCache(ctx, appCtx, updatedUserInfo.UserId); err != nil {
		logger.AppLogger.Error(
			ctx,
			"Failed updateUserCache",
			zap.Error(err),
		)
		msg.Ack()
		return
	}

	msg.Ack()
}
