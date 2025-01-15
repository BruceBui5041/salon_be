package couponbiz

import (
	"context"
	"errors"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/coupon/couponerror"
	"salon_be/model/coupon/couponmodel"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type CreateCouponRepo interface {
	CreateCoupon(ctx context.Context, data *couponmodel.CreateCoupon) error
}

type createCouponBiz struct {
	repo CreateCouponRepo
}

func NewCreateCouponBiz(repo CreateCouponRepo) *createCouponBiz {
	return &createCouponBiz{repo: repo}
}

func (biz *createCouponBiz) CreateCoupon(ctx context.Context, data *couponmodel.CreateCoupon) error {
	if data.EndDate.Before(data.StartDate) {
		err := errors.New("end date must be after start date")
		logger.AppLogger.Error(ctx, "Invalid coupon dates",
			zap.Error(err),
			zap.Time("start_date", data.StartDate),
			zap.Time("end_date", data.EndDate))
		return couponerror.ErrCouponInvalid(err)
	}

	if data.EndDate.Before(time.Now()) {
		err := errors.New("end date must be in the future")
		logger.AppLogger.Error(ctx, "Invalid coupon end date",
			zap.Error(err),
			zap.Time("end_date", data.EndDate))
		return couponerror.ErrCouponExpired(err)
	}

	if data.DiscountType != models.DiscountTypePercentage && data.DiscountType != models.DiscountTypeFixedPrice {
		err := errors.New("invalid discount type")
		logger.AppLogger.Error(ctx, "Invalid coupon discount type",
			zap.Error(err),
			zap.String("discount_type", string(data.DiscountType)))
		return couponerror.ErrCouponInvalid(err)
	}

	if data.DiscountType == models.DiscountTypePercentage && data.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		err := errors.New("percentage discount cannot exceed 100")
		logger.AppLogger.Error(ctx, "Invalid coupon percentage discount",
			zap.Error(err),
			zap.String("discount_value", data.DiscountValue.String()))
		return couponerror.ErrCouponInvalid(err)
	}

	if err := biz.repo.CreateCoupon(ctx, data); err != nil {
		logger.AppLogger.Error(ctx, "Failed to create coupon",
			zap.Error(err),
			zap.String("code", data.Code),
			zap.String("discount_type", string(data.DiscountType)),
			zap.String("discount_value", data.DiscountValue.String()))
		return couponerror.ErrCouponExists(err)
	}

	logger.AppLogger.Info(ctx, "Coupon created successfully",
		zap.String("code", data.Code),
		zap.String("discount_type", string(data.DiscountType)),
		zap.String("discount_value", data.DiscountValue.String()))

	return nil
}
