package couponrepo

import (
	"context"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/coupon/couponmodel"

	"go.uber.org/zap"
)

type UpdateCouponStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Coupon, error)
	Update(ctx context.Context, id uint32, data *models.Coupon) error
}

type UpdateImageRepo interface {
	CreateImageForCoupon(
		ctx context.Context,
		file *multipart.FileHeader,
		couponId uint32,
		userID uint32,
	) (*models.Image, error)
}

type updateCouponRepo struct {
	store     UpdateCouponStore
	imageRepo UpdateImageRepo
}

func NewUpdateCouponRepo(store UpdateCouponStore, imageRepo UpdateImageRepo) *updateCouponRepo {
	return &updateCouponRepo{
		store:     store,
		imageRepo: imageRepo,
	}
}

func (repo *updateCouponRepo) FindCoupon(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) (*models.Coupon, error) {
	return repo.store.FindOne(ctx, conditions, moreInfo...)
}

func (repo *updateCouponRepo) UpdateCoupon(ctx context.Context, id uint32, data *couponmodel.UpdateCoupon) error {
	var status string
	if data.Status == nil {
		status = common.StatusActive
	} else {
		status = common.StatusInactive
	}

	coupon := &models.Coupon{
		SQLModel:      common.SQLModel{Status: status},
		Description:   data.Description,
		DiscountType:  data.DiscountType,
		DiscountValue: data.DiscountValue,
		MinSpend:      data.MinSpend,
		MaxDiscount:   data.MaxDiscount,
		StartDate:     data.StartDate,
		EndDate:       data.EndDate,
		UsageLimit:    data.UsageLimit,
	}

	if err := repo.store.Update(ctx, id, coupon); err != nil {
		logger.AppLogger.Error(ctx, "Failed to update coupon",
			zap.Error(err),
			zap.Uint32("coupon_id", id),
			zap.Any("coupon_data", data))
		return err
	}

	if data.Image != nil {
		img, err := repo.imageRepo.CreateImageForCoupon(ctx, data.Image, id, data.CreatorID)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to create image for coupon",
				zap.Error(err),
				zap.Uint32("coupon_id", id),
				zap.Any("image_data", data.Image))
			return err
		}

		coupon.ImageID = &img.Id
		if err := repo.store.Update(ctx, id, coupon); err != nil {
			logger.AppLogger.Error(ctx, "Failed to update coupon with new image",
				zap.Error(err),
				zap.Uint32("coupon_id", id),
				zap.Uint32("image_id", img.Id))
			return err
		}
	}

	return nil
}
