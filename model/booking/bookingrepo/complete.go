package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"time"
)

type CompleteBookingStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Booking, error)
	Update(ctx context.Context, id uint32, data *models.Booking) error
}

type completeBookingRepo struct {
	bookingStore CompleteBookingStore
}

func NewCompleteBookingRepo(bookingStore CompleteBookingStore) *completeBookingRepo {
	return &completeBookingRepo{bookingStore: bookingStore}
}

func (repo *completeBookingRepo) CompleteBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CompleteBooking) error {
	booking, err := repo.bookingStore.FindOne(
		ctx,
		map[string]interface{}{"id": bookingId},
		"ServiceMan",
	)
	if err != nil {
		return common.ErrEntityNotFound(models.BookingEntityName, err)
	}

	// Check if the requester is the service man of the booking
	if booking.ServiceManID != data.UserID {
		return common.ErrNoPermission(errors.New("only service man of the booking can complete it"))
	}

	// Check if booking is in confirmed status
	if booking.BookingStatus != models.BookingStatusConfirmed {
		return common.ErrInvalidRequest(errors.New("only confirmed bookings can be completed"))
	}

	now := time.Now().UTC()
	booking.BookingStatus = models.BookingStatusCompleted
	booking.CompletedAt = &now

	if err := repo.bookingStore.Update(ctx, bookingId, booking); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
