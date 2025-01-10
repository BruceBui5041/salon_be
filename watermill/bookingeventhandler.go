package watermill

import (
	"context"
	"encoding/json"
	"fmt"
	"salon_be/component"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/notification/notificationbiz"
	"salon_be/model/notification/notificationrepo"
	"salon_be/model/notification/notificationstore"
	notificationdetailstore "salon_be/model/notificationdetails/notificationdetailsstore"
	"salon_be/watermill/messagemodel"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

const (
	NotificationTypeBooking = "booking"
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

	if err := handleBookingNotification(ctx, appCtx, bookingEventMsg); err != nil {
		logger.AppLogger.Error(ctx, "Error handling booking notification",
			zap.Error(err),
			zap.Any("bookingEvent", bookingEventMsg),
		)
		msg.Ack()
		return
	}

	msg.Ack()
}

func handleBookingNotification(
	ctx context.Context,
	appCtx component.AppContext,
	bookingEvent *messagemodel.BookingEventMsg,
) error {
	db := appCtx.GetMainDBConnection()

	// Get booking details with preloaded relationships
	bookingStore := bookingstore.NewSQLStore(db)
	booking, err := bookingStore.FindOne(
		ctx,
		map[string]interface{}{"id": bookingEvent.BookingID},
		"ServiceVersions",
		"ServiceVersions.Service",
		"User",
		"ServiceMan",
	)
	if err != nil {
		logger.AppLogger.Error(ctx, "Error finding booking", zap.Error(err))
		return fmt.Errorf("error finding booking: %w", err)
	}

	// Initialize notification stores and repositories
	notificationStore := notificationstore.NewSQLStore(db)
	notificationDetailStore := notificationdetailstore.NewSQLStore(db)
	repo := notificationrepo.NewCreateNotificationRepo(notificationStore, notificationDetailStore)
	biz := notificationbiz.NewCreateNotificationBiz(repo, repo)

	var recipientID uint32

	// Collect service slugs and service version IDs
	var serviceVersionIDs []uint32
	for _, sv := range booking.ServiceVersions {
		serviceVersionIDs = append(serviceVersionIDs, sv.Id)
	}

	var serviceIDs []uint32
	for _, sv := range booking.ServiceVersions {
		serviceIDs = append(serviceIDs, sv.ServiceID)
	}
	metadata := models.Metadata{
		"booking_id":          bookingEvent.BookingID,
		"event":               bookingEvent.Event,
		"service_version_ids": serviceVersionIDs,
		"service_ids":         serviceIDs,
	}

	recipientID = booking.UserID
	scheduled := time.Now()

	notification := &models.Notification{
		Type:      NotificationTypeBooking,
		BookingID: bookingEvent.BookingID,
		Metadata:  metadata,
		Scheduled: &scheduled,
		Details: []*models.NotificationDetail{
			{
				UserID: recipientID,
				State:  models.NotificationStatePending,
			},
		},
	}

	if err := biz.CreateNotification(ctx, notification); err != nil {
		logger.AppLogger.Error(ctx, "Error creating notification", zap.Error(err))
		return fmt.Errorf("error creating notification: %w", err)
	}

	return nil
}
