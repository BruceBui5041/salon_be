package bookingrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"

	"github.com/shopspring/decimal"
)

type BookingStore interface {
	Create(ctx context.Context, data *models.Booking) error
}

type ServiceVersionStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.ServiceVersion, error)
}

type ServiceManStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type createBookingRepo struct {
	bookingStore        BookingStore
	serviceVersionStore ServiceVersionStore
	serviceManStore     ServiceManStore
}

func NewCreateBookingRepo(
	bookingStore BookingStore,
	serviceVersionStore ServiceVersionStore,
	serviceManStore ServiceManStore,
) *createBookingRepo {
	return &createBookingRepo{
		bookingStore:        bookingStore,
		serviceVersionStore: serviceVersionStore,
		serviceManStore:     serviceManStore,
	}
}

func (repo *createBookingRepo) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) error {
	serviceVersion, err := repo.GetServiceVersion(ctx, data.ServiceVersionID)
	if err != nil {
		return err
	}

	_, err = repo.GetServiceMan(ctx, data.ServiceManID)
	if err != nil {
		return err
	}

	booking := &models.Booking{
		UserID:           data.UserID,
		ServiceVersionID: data.ServiceVersionID,
		ServiceManID:     data.ServiceManID,
		BookingDate:      data.BookingDate,
		Duration:         serviceVersion.Duration,
		Price:            serviceVersion.Price,
		Notes:            data.Notes,
		BookingStatus:    models.BookingStatusPending,
	}

	booking.DiscountAmount = decimal.NewFromInt(0)

	if err := repo.bookingStore.Create(ctx, booking); err != nil {
		return common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return nil
}

func (repo *createBookingRepo) GetServiceVersion(ctx context.Context, id uint32) (*models.ServiceVersion, error) {
	serviceVersion, err := repo.serviceVersionStore.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrEntityNotFound(models.ServiceVersionEntityName, err)
	}
	return serviceVersion, nil
}

func (repo *createBookingRepo) GetServiceMan(ctx context.Context, id uint32) (*models.User, error) {
	serviceMan, err := repo.serviceManStore.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}
	return serviceMan, nil
}
