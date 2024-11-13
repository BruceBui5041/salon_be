package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"time"
)

type AcceptBookingStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Booking, error)
	Update(ctx context.Context, id uint32, data *models.Booking) error
}

type acceptBookingRepo struct {
	store AcceptBookingStore
}

func NewAcceptBookingRepo(store AcceptBookingStore) *acceptBookingRepo {
	return &acceptBookingRepo{store: store}
}

func (repo *acceptBookingRepo) AcceptBooking(ctx context.Context, bookingId uint32, data *bookingmodel.AcceptBooking) error {
	booking, err := repo.store.FindOne(ctx, map[string]interface{}{"id": bookingId})
	if err != nil {
		return common.ErrEntityNotFound(models.BookingEntityName, err)
	}

	// Verify booking is in pending status
	if booking.BookingStatus != models.BookingStatusPending {
		return common.ErrInvalidRequest(errors.New("booking must be in pending status to be accepted"))
	}

	// Verify the user is the service provider of this booking
	if booking.ServiceManID != data.UserID {
		return common.ErrNoPermission(errors.New("only the assigned service provider can accept this booking"))
	}

	// Set confirmation time in UTC
	now := time.Now().UTC()
	booking.BookingStatus = models.BookingStatusConfirmed
	booking.ConfirmedDate = &now

	if err := repo.store.Update(ctx, bookingId, booking); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
