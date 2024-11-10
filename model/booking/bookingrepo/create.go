package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"

	"github.com/shopspring/decimal"
)

type BookingStore interface {
	Create(ctx context.Context, data *models.Booking) error
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Booking, error)
}

type ServiceStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Service, error)
}

type createBookingRepo struct {
	bookingStore BookingStore
	serviceStore ServiceStore
}

func NewCreateBookingRepo(
	bookingStore BookingStore,
	serviceStore ServiceStore,
) *createBookingRepo {
	return &createBookingRepo{
		bookingStore: bookingStore,
		serviceStore: serviceStore,
	}
}

func (repo *createBookingRepo) checkUserPendingBookings(ctx context.Context, userID uint32) error {
	conditions := map[string]interface{}{
		"user_id": userID,
	}

	bookings, err := repo.bookingStore.Find(ctx, conditions)
	if err != nil {
		return common.ErrCannotListEntity(models.BookingEntityName, err)
	}

	for _, booking := range bookings {
		if booking.BookingStatus != models.BookingStatusCompleted &&
			booking.BookingStatus != models.BookingStatusCancelled {
			return common.ErrInvalidRequest(errors.New("cannot create new booking while having pending bookings"))
		}
	}

	return nil
}

func (repo *createBookingRepo) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) error {
	if !data.IsUserRole {
		return common.ErrNoPermission(errors.New("only users can create bookings"))
	}

	if err := repo.checkUserPendingBookings(ctx, data.UserID); err != nil {
		return err
	}

	// Get service with its current version and creator
	serviceID, err := data.GetVersionLocalId(ctx)
	if err != nil {
		return err
	}

	service, err := repo.serviceStore.FindOne(
		ctx,
		map[string]interface{}{"id": serviceID},
		"ServiceVersion",
		"Creator",
	)
	if err != nil {
		return common.ErrEntityNotFound(models.ServiceEntityName, err)
	}

	if service.ServiceVersion == nil {
		return common.ErrEntityNotFound(models.ServiceVersionEntityName, errors.New("service version not found"))
	}

	if service.Creator == nil {
		return common.ErrEntityNotFound(models.UserEntityName, errors.New("service creator not found"))
	}

	booking := &models.Booking{
		UserID:           data.UserID,
		ServiceVersionID: service.ServiceVersion.Id,
		ServiceManID:     service.Creator.Id,
		BookingDate:      data.BookingDate,
		Price:            service.ServiceVersion.Price,
		Notes:            data.Notes,
		ServiceVersion:   service.ServiceVersion,
	}

	booking.DiscountAmount = decimal.NewFromInt(0)

	if err := repo.bookingStore.Create(ctx, booking); err != nil {
		return common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return nil
}
