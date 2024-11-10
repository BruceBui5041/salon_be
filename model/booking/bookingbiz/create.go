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
	GetServiceVersion(ctx context.Context, id uint32) (*models.ServiceVersion, error)
	GetServiceMan(ctx context.Context, id uint32) (*models.User, error)
}

type createBookingBiz struct {
	repo BookingRepo
}

func NewCreateBookingBiz(repo BookingRepo) *createBookingBiz {
	return &createBookingBiz{repo: repo}
}

func (biz *createBookingBiz) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) error {
	if data.ServiceVersionID == 0 {
		return common.ErrInvalidRequest(errors.New("service version ID is required"))
	}

	if data.ServiceManID == 0 {
		return common.ErrInvalidRequest(errors.New("service man ID is required"))
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
