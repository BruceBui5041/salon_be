package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/payment/paymentconst"
	"time"
)

type CancelBookingStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Booking, error)
	Update(ctx context.Context, id uint32, data *models.Booking) error
}

type CancelBookingPaymentStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Payment, error)
	Update(ctx context.Context, id uint32, data *models.Payment) error
}

type cancelBookingRepo struct {
	store        CancelBookingStore
	paymentStore CancelBookingPaymentStore
}

func NewCancelBookingRepo(store CancelBookingStore, paymentStore CancelBookingPaymentStore) *cancelBookingRepo {
	return &cancelBookingRepo{
		store:        store,
		paymentStore: paymentStore,
	}
}

func (repo *cancelBookingRepo) CancelBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CancelBooking) error {
	booking, err := repo.store.FindOne(ctx, map[string]interface{}{"id": bookingId}, "Payment")
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
	if booking.UserID != data.UserID && booking.ServiceManID != data.UserID {
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

	if booking.PaymentID != nil {
		payment := &models.Payment{
			TransactionStatus: paymentconst.TransactionStatusCancelled,
		}

		if err := repo.paymentStore.Update(ctx, *booking.PaymentID, payment); err != nil {
			return common.ErrCannotUpdateEntity(models.PaymentEntityName, err)
		}
	}

	return nil
}
