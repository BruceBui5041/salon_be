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

type CouponStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Coupon, error)
}

type createBookingRepo struct {
	bookingStore BookingStore
	serviceStore ServiceStore
	paymentStore PaymentStore
	couponStore  CouponStore
}

func NewCreateBookingRepo(
	bookingStore BookingStore,
	serviceStore ServiceStore,
	paymentStore PaymentStore,
	couponStore CouponStore,
) *createBookingRepo {
	return &createBookingRepo{
		bookingStore: bookingStore,
		serviceStore: serviceStore,
		paymentStore: paymentStore,
		couponStore:  couponStore,
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

	// Handle coupon if provided
	var coupon *models.Coupon
	if data.CouponID != nil {
		couponID, err := data.GetCouponLocalId()
		if err != nil {
			return 0, err
		}

		coupon, err = repo.couponStore.FindOne(ctx, map[string]interface{}{
			"id":     couponID,
			"code":   data.CouponCode,
			"status": common.StatusActive,
		})
		if err != nil {
			if err == common.RecordNotFound {
				return 0, common.ErrEntityNotFound(models.CouponEntityName, errors.New("coupon not found"))
			}
			return 0, common.ErrCannotGetEntity(models.CouponEntityName, err)
		}

		// Validate coupon
		if err := coupon.IsValid(totalPrice); err != nil {
			logger.AppLogger.Error(ctx, "invalid coupon", zap.Error(err), zap.Any("coupon", coupon))
			return 0, err
		}
	}

	// Create booking entity with UTC time
	booking := &models.Booking{
		UserID:          data.UserID,
		ServiceVersions: serviceVersions,
		ServiceManID:    serviceManID,
		BookingDate:     data.BookingDate.UTC(),
		Duration:        totalDuration,
		Price:           totalPrice,
		Notes:           data.Notes,
	}

	// Set coupon if provided
	if coupon != nil {
		booking.CouponID = &coupon.Id
		booking.Coupon = coupon
	}

	// Calculate discounted price using the model's method
	if err := booking.CalculateDiscountedPrice(); err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	// Create payment record with final price
	finalPrice := totalPrice
	if booking.DiscountedPrice != nil {
		finalPrice = *booking.DiscountedPrice
	}

	payment, err := repo.createPayment(ctx, data.UserID, finalPrice, data.PaymentMethod)
	if err != nil {
		logger.AppLogger.Error(ctx, "cannot create payment", zap.Error(err), zap.Any("data", data))
		return 0, err
	}

	booking.PaymentID = &payment.Id
	booking.Payment = payment

	id, err := repo.bookingStore.Create(ctx, booking)
	if err != nil {
		return 0, common.ErrCannotCreateEntity(models.BookingEntityName, err)
	}

	return id, nil
}
