// booking/bookingrepo/cancel.go
package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"time"
)

type CancelBookingStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Booking, error)
	Update(ctx context.Context, id uint32, data *models.Booking) error
}

type cancelBookingRepo struct {
	store CancelBookingStore
}

func NewCancelBookingRepo(store CancelBookingStore) *cancelBookingRepo {
	return &cancelBookingRepo{store: store}
}

func (repo *cancelBookingRepo) CancelBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CancelBooking) error {
	booking, err := repo.store.FindOne(ctx, map[string]interface{}{"id": bookingId})
	if err != nil {
		return common.ErrEntityNotFound(models.BookingEntityName, err)
	}

	if booking.BookingStatus == models.BookingStatusCancelled {
		return common.ErrInvalidRequest(errors.New("booking is already cancelled"))
	}

	if booking.BookingStatus == models.BookingStatusCompleted {
		return common.ErrInvalidRequest(errors.New("cannot cancel completed booking"))
	}

	// Check if the requester is either the service man or the booking user
	isAuthorized := false
	if data.IsUserRole && booking.UserID == data.UserID {
		isAuthorized = true
	} else if !data.IsUserRole && booking.ServiceManID == data.UserID {
		isAuthorized = true
	}

	if !isAuthorized {
		return common.ErrNoPermission(errors.New("only the service provider or the customer can cancel this booking"))
	}

	// Set cancellation time in UTC
	now := time.Now().UTC()
	booking.BookingStatus = models.BookingStatusCancelled
	booking.CancellationReason = data.CancellationReason
	booking.CancelledByID = &data.UserID
	booking.CancelledAt = &now

	if err := repo.store.Update(ctx, bookingId, booking); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
