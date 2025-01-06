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

	// Get booking details
	bookingStore := bookingstore.NewSQLStore(db)
	booking, err := bookingStore.FindOne(
		ctx,
		map[string]interface{}{"id": bookingEvent.BookingID},
		"ServiceVersion",
		"ServiceVersion.Service",
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

	// Create notification based on event type
	var title, content string
	var recipientID uint32
	metadata := models.Metadata{
		"booking_id": bookingEvent.BookingID,
		"event":      bookingEvent.Event,
	}

	recipientID = booking.UserID
	scheduled := time.Now()
	switch bookingEvent.Event {
	case messagemodel.BookingCreatedEvent:
		title = "New Booking Request"
		content = fmt.Sprintf("New booking request for service: %s", booking.ServiceVersion.Service.Slug)
		metadata["service_id"] = booking.ServiceVersion.ServiceID

	case messagemodel.BookingAcceptedEvent:
		title = "Booking Accepted"
		content = fmt.Sprintf("Your booking for %s has been accepted", booking.ServiceVersion.Service.Slug)

	case messagemodel.BookingCompletedEvent:
		title = "Booking Completed"
		content = fmt.Sprintf("Your booking for %s has been completed", booking.ServiceVersion.Service.Slug)

	case messagemodel.BookingCancelledEvent:
		if *booking.CancelledByID == booking.UserID {
			title = "Booking Cancelled by Customer"
			content = fmt.Sprintf("Booking for %s has been cancelled by the customer", booking.ServiceVersion.Service.Slug)
		} else {
			title = "Booking Cancelled by Service Provider"
			content = fmt.Sprintf("Your booking for %s has been cancelled by the service provider", booking.ServiceVersion.Service.Slug)
		}

	default:
		logger.AppLogger.Error(ctx, "Unknown booking event", zap.Any("bookingEvent", bookingEvent))
		return fmt.Errorf("unknown booking event: %s", bookingEvent.Event)
	}

	notification := &models.Notification{
		Title:     title,
		Content:   content,
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
