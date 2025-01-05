package bookingbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/payment/paymentconst"
	"time"

	"github.com/samber/lo"
)

type BookingRepo interface {
	CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) (uint32, error)
}

type createBookingBiz struct {
	repo BookingRepo
}

func NewCreateBookingBiz(repo BookingRepo) *createBookingBiz {
	return &createBookingBiz{repo: repo}
}

func (biz *createBookingBiz) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) (uint32, error) {
	if data.ServiceID == "" {
		return 0, common.ErrInvalidRequest(errors.New("service ID is required"))
	}

	if data.BookingDate.IsZero() {
		return 0, common.ErrInvalidRequest(errors.New("booking date is required"))
	}

	if data.PaymentMethod == "" {
		return 0, common.ErrInvalidRequest(errors.New("payment method is required"))
	}

	if lo.IndexOf(paymentconst.PaymentMethods, data.PaymentMethod) == -1 {
		return 0, common.ErrInvalidRequest(errors.New("invalid payment method"))
	}

	// Check if booking date is in the future
	if data.BookingDate.Before(time.Now().UTC()) {
		return 0, common.ErrInvalidRequest(errors.New("booking date must be in the future"))
	}

	id, err := biz.repo.CreateBooking(ctx, data)
	if err != nil {
		return 0, common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return id, nil
}
