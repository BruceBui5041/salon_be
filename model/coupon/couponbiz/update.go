package couponbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/coupon/couponerror"
	"salon_be/model/coupon/couponmodel"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UpdateCouponRepo interface {
	FindCoupon(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Coupon, error)
	UpdateCoupon(ctx context.Context, id uint32, data *couponmodel.UpdateCoupon) error
}

type updateCouponBiz struct {
	repo UpdateCouponRepo
}

func NewUpdateCouponBiz(repo UpdateCouponRepo) *updateCouponBiz {
	return &updateCouponBiz{repo: repo}
}

func (biz *updateCouponBiz) UpdateCoupon(ctx context.Context, id string, data *couponmodel.UpdateCoupon) error {
	uid, err := common.FromBase58(id)
	if err != nil {
		return couponerror.ErrCouponInvalid(err)
	}

	// Find existing coupon
	coupon, err := biz.repo.FindCoupon(ctx, map[string]interface{}{"id": uid.GetLocalID()})
	if err != nil {
		return common.ErrCannotGetEntity(models.CouponEntityName, err)
	}

	// Check if coupon is inactive
	if coupon.Status != common.StatusInactive {
		return couponerror.ErrCouponInvalid(errors.New("can only update inactive coupons"))
	}

	// Check date conditions
	nowUTC := time.Now().UTC()
	if nowUTC.Before(coupon.StartDate) || nowUTC.After(coupon.EndDate) {
		return couponerror.ErrCouponInvalid(errors.New("coupon can only be updated within its valid date range"))
	}

	if data.EndDate.Before(data.StartDate) {
		err := errors.New("end date must be after start date")
		logger.AppLogger.Error(ctx, "Invalid coupon dates",
			zap.Error(err),
			zap.Time("start_date", data.StartDate),
			zap.Time("end_date", data.EndDate))
		return couponerror.ErrCouponInvalid(err)
	}

	if data.DiscountType == models.DiscountTypePercentage && data.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		err := errors.New("percentage discount cannot exceed 100")
		logger.AppLogger.Error(ctx, "Invalid coupon percentage discount",
			zap.Error(err),
			zap.String("discount_value", data.DiscountValue.String()))
		return couponerror.ErrCouponInvalid(err)
	}

	if err := biz.repo.UpdateCoupon(ctx, uid.GetLocalID(), data); err != nil {
		return common.ErrCannotUpdateEntity(models.CouponEntityName, err)
	}

	return nil
}
