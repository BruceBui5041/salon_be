package bookingrepo

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/payment/paymentconst"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BookingStore interface {
	Create(ctx context.Context, data *models.Booking) (uint32, error)
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Booking, error)
}

type ServiceStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Service, error)
	Find(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) ([]*models.Service, error)
}

type PaymentStore interface {
	Create(ctx context.Context, data *models.Payment) (uint32, error)
}

type createBookingRepo struct {
	bookingStore BookingStore
	serviceStore ServiceStore
	paymentStore PaymentStore
}

func NewCreateBookingRepo(
	bookingStore BookingStore,
	serviceStore ServiceStore,
	paymentStore PaymentStore,
) *createBookingRepo {
	return &createBookingRepo{
		bookingStore: bookingStore,
		serviceStore: serviceStore,
		paymentStore: paymentStore,
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

func (repo *createBookingRepo) createPayment(
	ctx context.Context,
	userID uint32,
	amount decimal.Decimal,
	paymentMethod string,
) (*models.Payment, error) {
	payment := &models.Payment{
		UserID:            userID,
		Amount:            amount.InexactFloat64(),
		Currency:          "VND",
		PaymentMethod:     paymentMethod,
		TransactionStatus: paymentconst.TransactionStatusPending,
		TransactionID:     uuid.New().String(),
	}

	paymentID, err := repo.paymentStore.Create(ctx, payment)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.PaymentEntityName, err)
	}

	payment.Id = paymentID
	return payment, nil
}

func (repo *createBookingRepo) CreateBooking(ctx context.Context, data *bookingmodel.CreateBooking) (uint32, error) {
	if !data.IsUserRole {
		return 0, common.ErrNoPermission(errors.New("only users can create bookings"))
	}

	if err := repo.checkUserPendingBookings(ctx, data.UserID); err != nil {
		logger.AppLogger.Error(ctx, "cannot create booking", zap.Error(err), zap.Any("data", data))
		return 0, err
	}

	serviceIds, err := data.GetVersionLocalIds(ctx)
	if err != nil {
		return 0, err
	}

	services, err := repo.serviceStore.Find(
		ctx,
		map[string]interface{}{
			"id": serviceIds,
		},
		"ServiceVersion",
		"Creator",
	)
	if err != nil {
		return 0, common.ErrEntityNotFound(models.ServiceEntityName, err)
	}

	if len(services) != len(serviceIds) {
		logger.AppLogger.Error(
			ctx,
			"some services not found",
			zap.Error(err),
			zap.Any("service_ids", serviceIds),
			zap.Any("found services", services),
		)
		return 0, common.ErrEntityNotFound(models.ServiceEntityName, errors.New("some services not found"))
	}

	var serviceVersions []*models.ServiceVersion
	var totalDuration uint32
	var totalPrice decimal.Decimal
	var serviceManID uint32

	// Process all services
	for _, service := range services {
		if service.ServiceVersion == nil {
			logger.AppLogger.Error(ctx, "service version not found", zap.Any("service", service))
			return 0, common.ErrEntityNotFound(models.ServiceVersionEntityName, errors.New("service version not found"))
		}

		if service.Creator == nil {
			logger.AppLogger.Error(ctx, "service creator not found", zap.Any("service", service))
			return 0, common.ErrEntityNotFound(models.UserEntityName, errors.New("service creator not found"))
		}

		// Use the first service's creator as the service man
		if serviceManID == 0 {
			serviceManID = service.Creator.Id
		}

		serviceVersions = append(serviceVersions, service.ServiceVersion)
		totalDuration += service.ServiceVersion.Duration

		// Add to total price
		if service.ServiceVersion.DiscountedPrice != nil {
			totalPrice = totalPrice.Add(service.ServiceVersion.DiscountedPrice.Decimal)
		} else {
			totalPrice = totalPrice.Add(service.ServiceVersion.Price)
		}
	}

	// Create payment record
	payment, err := repo.createPayment(ctx, data.UserID, totalPrice, data.PaymentMethod)
	if err != nil {
		logger.AppLogger.Error(ctx, "cannot create payment", zap.Error(err), zap.Any("data", data))
		return 0, err
	}

	booking := &models.Booking{
		UserID:          data.UserID,
		ServiceVersions: serviceVersions,
		ServiceManID:    serviceManID,
		BookingDate:     data.BookingDate,
		Duration:        totalDuration,
		Price:           totalPrice,
		Notes:           data.Notes,
		PaymentID:       &payment.Id,
		Payment:         payment,
	}

	booking.DiscountAmount = decimal.NewFromInt(0)

	id, err := repo.bookingStore.Create(ctx, booking)
	if err != nil {
		return 0, common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return id, nil
}
