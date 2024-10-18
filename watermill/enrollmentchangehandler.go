package watermill

import (
	"encoding/json"
	"video_server/component"
	"video_server/component/logger"
	"video_server/model/enrollment/enrollmentutils"
	"video_server/model/user/userutils"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func EnrollmentChangeHandler(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "EnrollmentChangeHandler")
	defer span.End()

	var updateUserCacheInfo *messagemodel.EnrollmentChangeInfo
	err := json.Unmarshal(msg.Payload, &updateUserCacheInfo)
	if err != nil {
		logger.AppLogger.Error(ctx, "Cannot unmarshal message payload", zap.Any("payload", msg.Payload), zap.Error(err))
		msg.Ack()
		return
	}

	if err := userutils.UpdateUserCache(ctx, appCtx, updateUserCacheInfo.UserId); err != nil {
		logger.AppLogger.Error(ctx,
			"Failed updateUserCache",
			zap.Error(err),
		)

		msg.Ack()
		return
	}

	if err := enrollmentutils.UpdateEnrollmentCache(ctx, appCtx, updateUserCacheInfo); err != nil {
		logger.AppLogger.Error(ctx,
			"Failed updateEnrollmentCache",
			zap.Error(err),
		)

		msg.Ack()
		return
	}

	msg.Ack()
}
