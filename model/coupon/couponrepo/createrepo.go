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

type CreateCouponStore interface {
	Create(ctx context.Context, data *models.Coupon) error
	Update(ctx context.Context, id uint32, data *models.Coupon) error
}

type CreateImageRepo interface {
	CreateImage(ctx context.Context, file *multipart.FileHeader, serviceID uint32, userID uint32) (*models.Image, error)
}

type createCouponRepo struct {
	store     CreateCouponStore
	imageRepo CreateImageRepo
}

func NewCreateCouponRepo(store CreateCouponStore, imageRepo CreateImageRepo) *createCouponRepo {
	return &createCouponRepo{
		store:     store,
		imageRepo: imageRepo,
	}
}

func (repo *createCouponRepo) CreateCoupon(ctx context.Context, data *couponmodel.CreateCoupon) (uint32, error) {
	coupon := &models.Coupon{
		SQLModel:      common.SQLModel{Status: common.StatusInactive},
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
		logger.AppLogger.Error(ctx, "Failed to create coupon in database",
			zap.Error(err),
			zap.String("code", data.Code))
		return 0, err
	}

	if data.Image != nil {
		img, err := repo.imageRepo.CreateImage(ctx, data.Image, coupon.Id, data.CreatorID)
		if err != nil {
			logger.AppLogger.Error(ctx, "failed to upload coupon image", zap.Error(err))
			return 0, err
		}

		coupon.ImageID = &img.Id
		if err := repo.store.Update(ctx, coupon.Id, coupon); err != nil {
			logger.AppLogger.Error(ctx, "failed to update coupon with image", zap.Error(err))
			return 0, err
		}
	}

	return coupon.Id, nil
}
