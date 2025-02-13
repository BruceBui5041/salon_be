package watermill

import (
	"encoding/json"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/location/locationbiz"
	"salon_be/model/location/locationrepo"
	"salon_be/model/location/locationstore"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func ReceivedWSMsgHandler(appCtx component.AppContext, msg *message.Message) {
	ctx, span := createTracedHandler(msg, "ReceivedWSMsgHandler")
	defer span.End()

	logger.AppLogger.Info(ctx, "msg userID", zap.Any("user id", msg.Metadata.Get("currentUserID")))
	logger.AppLogger.Info(ctx, "ReceivedWSMsgHandler", zap.Any("msg payload", msg.Payload))

	var wsMessage messagemodel.UpdateLocationBody
	err := json.Unmarshal(msg.Payload, &wsMessage)
	if err != nil {
		msg.Ack()
		logger.AppLogger.Error(ctx, "Cannot unmarshal message payload",
			zap.Any("payload", msg.Payload),
			zap.Error(err))
		return
	}

	// Handle only location update events
	if wsMessage.Event != "location_update" {
		msg.Ack()
		return
	}

	wsMessage.UserId = msg.Metadata.Get("currentUserID")

	// Convert string ID to uint32
	uid, err := common.FromBase58(wsMessage.UserId)
	if err != nil {
		msg.Ack()
		logger.AppLogger.Error(ctx, "Invalid user ID format",
			zap.String("user_id", wsMessage.UserId),
			zap.Error(err))
		return
	}

	// Initialize the location update components
	db := appCtx.GetMainDBConnection()
	locationStore := locationstore.NewSQLStore(db)
	repo := locationrepo.NewUpdateLocationRepo(locationStore)
	biz := locationbiz.NewUpdateLocationBiz(repo)

	// Create location data
	location := &models.Location{
		UserId:    uid.GetLocalID(),
		Latitude:  wsMessage.Latitude,
		Longitude: wsMessage.Longitude,
		Accuracy:  wsMessage.Accuracy,
	}

	// Update location
	if err := biz.UpdateLocation(ctx, location); err != nil {
		msg.Ack()
		logger.AppLogger.Error(ctx, "Failed to update location",
			zap.Any("location", location),
			zap.Error(err))
		return
	}

	logger.AppLogger.Info(ctx, "Location updated successfully",
		zap.String("user_id", wsMessage.UserId))

	msg.Ack()
}
