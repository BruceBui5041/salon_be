package couponbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/coupon/couponmodel"
	"time"

	"github.com/shopspring/decimal"
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
		return common.ErrInvalidRequest(errors.New("end date must be after start date"))
	}

	if data.EndDate.Before(time.Now()) {
		return common.ErrInvalidRequest(errors.New("end date must be in the future"))
	}

	if data.DiscountType != models.DiscountTypePercentage && data.DiscountType != models.DiscountTypeFixedPrice {
		return common.ErrInvalidRequest(errors.New("invalid discount type"))
	}

	if data.DiscountType == models.DiscountTypePercentage && data.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		return common.ErrInvalidRequest(errors.New("percentage discount cannot exceed 100"))
	}

	if err := biz.repo.CreateCoupon(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(models.CouponEntityName, err)
	}

	return nil
}
