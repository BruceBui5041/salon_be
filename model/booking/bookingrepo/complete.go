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

type CompleteBookingStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Booking, error)
	Update(ctx context.Context, id uint32, data *models.Booking) error
}

type CompleteBookingPaymentStore interface {
	Update(ctx context.Context, id uint32, data *models.Payment) error
}

type completeBookingRepo struct {
	bookingStore CompleteBookingStore
	paymentStore CompleteBookingPaymentStore
}

func NewCompleteBookingRepo(bookingStore CompleteBookingStore, paymentStore CompleteBookingPaymentStore) *completeBookingRepo {
	return &completeBookingRepo{
		bookingStore: bookingStore,
		paymentStore: paymentStore,
	}
}

func (repo *completeBookingRepo) CompleteBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CompleteBooking) error {
	booking, err := repo.bookingStore.FindOne(
		ctx,
		map[string]interface{}{"id": bookingId},
		"ServiceMan",
		"Payment",
	)
	if err != nil {
		return common.ErrEntityNotFound(models.BookingEntityName, err)
	}

	if booking.ServiceManID != data.UserID {
		return common.ErrNoPermission(errors.New("only service man of the booking can complete it"))
	}

	if booking.BookingStatus != models.BookingStatusConfirmed {
		return common.ErrInvalidRequest(errors.New("only accepted bookings can be completed"))
	}

	now := time.Now().UTC()
	booking.BookingStatus = models.BookingStatusCompleted
	booking.CompletedAt = &now

	if booking.Payment != nil {
		updatePayment := &models.Payment{TransactionStatus: paymentconst.TransactionStatusCompleted}
		if err := repo.paymentStore.Update(ctx, *booking.PaymentID, updatePayment); err != nil {
			return common.ErrCannotUpdateEntity(models.PaymentEntityName, err)
		}
	}

	if err := repo.bookingStore.Update(ctx, bookingId, booking); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
