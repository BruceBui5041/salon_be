package bookingbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
)

type CompleteBookingRepo interface {
	CompleteBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CompleteBooking) error
}

type completeBookingBiz struct {
	repo CompleteBookingRepo
}

func NewCompleteBookingBiz(repo CompleteBookingRepo) *completeBookingBiz {
	return &completeBookingBiz{repo: repo}
}

func (biz *completeBookingBiz) CompleteBooking(ctx context.Context, bookingId uint32, data *bookingmodel.CompleteBooking) error {
	if bookingId == 0 {
		return common.ErrInvalidRequest(errors.New("booking ID is required"))
	}

	if err := biz.repo.CompleteBooking(ctx, bookingId, data); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
