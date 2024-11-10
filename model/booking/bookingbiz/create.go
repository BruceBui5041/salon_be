// booking/bookingbiz/create.go

package bookingbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"time"
)

type BookingRepo interface {
	CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) error
}

type createBookingBiz struct {
	repo BookingRepo
}

func NewCreateBookingBiz(repo BookingRepo) *createBookingBiz {
	return &createBookingBiz{repo: repo}
}

func (biz *createBookingBiz) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) error {
	if data.ServiceID == "" {
		return common.ErrInvalidRequest(errors.New("service ID is required"))
	}

	if data.BookingDate.IsZero() {
		return common.ErrInvalidRequest(errors.New("booking date is required"))
	}

	// Check if booking date is in the future
	if data.BookingDate.Before(time.Now()) {
		return common.ErrInvalidRequest(errors.New("booking date must be in the future"))
	}

	if err := biz.repo.CreateBooking(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return nil
}
