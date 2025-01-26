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
	ServiceIDs    []string  `json:"service_ids"` // Changed from ServiceID to ServiceIDs
	CouponID      *string   `json:"coupon_id"`
	CouponCode    *string   `json:"coupon_code"`
	BookingDate   time.Time `json:"booking_date"`
	Notes         string    `json:"notes"`
	PaymentMethod string    `json:"payment_method"`
	UserID        uint32    `json:"-"`
	IsUserRole    bool      `json:"-"`
}

// Add new method to get multiple version local IDs
func (cb *CreateBooking) GetCouponLocalId() (uint32, error) {
	if cb.CouponID == nil {
		return 0, nil
	}

	couponUID, err := common.FromBase58(*cb.CouponID)
	if err != nil {
		return 0, common.ErrInvalidRequest(errors.New("invalid coupon ID"))
	}
	return couponUID.GetLocalID(), nil
}

func (cb *CreateBooking) GetVersionLocalIds(ctx context.Context) ([]uint32, error) {
	var localIds []uint32
	for _, serviceID := range cb.ServiceIDs {
		serviceUID, err := common.FromBase58(serviceID)
		if err != nil {
			logger.AppLogger.Error(ctx, "invalid service ID", zap.Error(err))
			return nil, common.ErrInvalidRequest(errors.New("invalid service ID"))
		}
		localIds = append(localIds, serviceUID.GetLocalID())
	}
	return localIds, nil
}
