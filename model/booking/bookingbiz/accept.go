package bookingbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
)

type AcceptBookingRepo interface {
	AcceptBooking(ctx context.Context, bookingId uint32, data *bookingmodel.AcceptBooking) error
}

type acceptBookingBiz struct {
	repo AcceptBookingRepo
}

func NewAcceptBookingBiz(repo AcceptBookingRepo) *acceptBookingBiz {
	return &acceptBookingBiz{repo: repo}
}

func (biz *acceptBookingBiz) AcceptBooking(ctx context.Context, bookingId uint32, data *bookingmodel.AcceptBooking) error {
	if bookingId == 0 {
		return common.ErrInvalidRequest(errors.New("booking ID is required"))
	}

	if err := biz.repo.AcceptBooking(ctx, bookingId, data); err != nil {
		return common.ErrCannotUpdateEntity(models.BookingEntityName, err)
	}

	return nil
}
