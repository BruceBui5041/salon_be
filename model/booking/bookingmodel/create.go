package bookingmodel

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component/logger"
	"time"

	"go.uber.org/zap"
)

type CreateBooking struct {
	ServiceID   string    `json:"service_id"`
	CouponID    *string   `json:"coupon_id"`
	BookingDate time.Time `json:"booking_date"`
	Notes       string    `json:"notes"`
	UserID      uint32    `json:"-"`
	IsUserRole  bool      `json:"-"`
}

func (cb *CreateBooking) GetVersionLocalId(ctx context.Context) (uint32, error) {
	serviceUID, err := common.FromBase58(cb.ServiceID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid service ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid service ID"))
	}
	return serviceUID.GetLocalID(), nil
}

func (cb *CreateBooking) GetCouponLocalId(ctx context.Context) (uint32, error) {
	if cb.CouponID == nil {
		return 0, nil
	}

	couponUID, err := common.FromBase58(*cb.CouponID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid coupon ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid coupon ID"))
	}
	return couponUID.GetLocalID(), nil
}
