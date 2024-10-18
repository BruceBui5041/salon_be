package watermill

import (
	"encoding/json"
	"salon_be/component"
	"salon_be/component/appqueue/providerhandler"
	"salon_be/component/logger"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func HandleNewVideoUpload(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "HandleNewVideoUpload")
	defer span.End()

	var videoInfo *messagemodel.RequestProcessVideoInfo
	err := json.Unmarshal(msg.Payload, &videoInfo)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"Cannot unmarshal message payload",
			zap.Any("payload", msg.Payload),
			zap.Error(err),
		)
		return
	}

	err = providerhandler.SendRequestProcessVideo(ctx, appCtx.GetAppQueue(), videoInfo)
	if err != nil {
		logger.AppLogger.Error(
			ctx,
			"SendRequestProcessVideo err",
			zap.Any("videoInfo", videoInfo),
			zap.Error(err),
		)
		return
	}

	// go grpcserver.ProcessNewVideoRequest(appCtx, videoInfo)
	msg.Ack()
}
