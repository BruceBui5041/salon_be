package watermill

import (
	"encoding/json"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func BookingEventHandler(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "BookingEventHandler")
	defer span.End()

	var bookingEventMsg *messagemodel.BookingEventMsg
	err := json.Unmarshal(msg.Payload, &bookingEventMsg)
	if err != nil {
		logger.AppLogger.Error(ctx, "Cannot unmarshal message payload", zap.Any("payload", msg.Payload), zap.Error(err))
		msg.Ack()
		return
	}

	logger.AppLogger.Info(ctx, "Received booking event", zap.Any("bookingEvent", bookingEventMsg))
}
