package couponrepo

import (
	"context"
	models "salon_be/model"
	"salon_be/model/coupon/couponmodel"
)

type CreateCouponStore interface {
	Create(ctx context.Context, data *models.Coupon) error
}

type createCouponRepo struct {
	store CreateCouponStore
}

func NewCreateCouponRepo(store CreateCouponStore) *createCouponRepo {
	return &createCouponRepo{store: store}
}

func (repo *createCouponRepo) CreateCoupon(ctx context.Context, data *couponmodel.CreateCoupon) error {
	coupon := &models.Coupon{
		Code:          data.Code,
		Description:   data.Description,
		DiscountType:  data.DiscountType,
		DiscountValue: data.DiscountValue,
		MinSpend:      data.MinSpend,
		MaxDiscount:   data.MaxDiscount,
		StartDate:     data.StartDate,
		EndDate:       data.EndDate,
		UsageLimit:    data.UsageLimit,
		UsageCount:    0,
		CreatorID:     data.CreatorID,
	}

	if err := repo.store.Create(ctx, coupon); err != nil {
		return err
	}

	return nil
}
