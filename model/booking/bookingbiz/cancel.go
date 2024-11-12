package bookingbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
)

type CancelBookingRepo interface {
	CancelBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CancelBooking) error
}

type cancelBookingBiz struct {
	repo CancelBookingRepo
}

func NewCancelBookingBiz(repo CancelBookingRepo) *cancelBookingBiz {
	return &cancelBookingBiz{repo: repo}
}

func (biz *cancelBookingBiz) CancelBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CancelBooking) error {
	if bookingId == 0 {
		return common.ErrInvalidRequest(errors.New("booking ID is required"))
	}

	if data.CancellationReason == "" {
		return common.ErrInvalidRequest(errors.New("cancellation reason is required"))
	}

	if err := biz.repo.CancelBooking(ctx, bookingId, data); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
