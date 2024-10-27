package servicemodel

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	"salon_be/utils/customtypes"

	"go.uber.org/zap"
)

type UpdateServiceRequest struct {
	JSON   string                  `form:"json"`
	Images []*multipart.FileHeader `form:"images"`
}

type UpdateService struct {
	ServiceID string `json:"id" form:"id"`
	// ServiceVersionID string                `json:"service_version_id" form:"service_version_id"`
	ServiceVersion *UpdateServiceVersion `json:"service_version,omitempty" form:"service_version,omitempty"`
}

type UpdateServiceVersion struct {
	ID              string                         `json:"id" form:"id"`
	Title           string                         `json:"title" form:"title"`
	Description     string                         `json:"description" form:"description"`
	CategoryID      string                         `json:"category_id" form:"category_id"`
	SubCategoryID   string                         `json:"sub_category_id" form:"sub_category_id"`
	IntroVideoID    string                         `json:"intro_video_id,omitempty" form:"intro_video_id,omitempty"`
	Thumbnail       string                         `json:"thumbnail" form:"thumbnail"`
	Price           customtypes.DecimalString      `json:"price" form:"price"`
	DiscountedPrice *customtypes.NullDecimalString `json:"discounted_price,omitempty" form:"discounted_price,omitempty"`
	Duration        uint32                         `json:"duration" form:"duration"`
	Images          []*multipart.FileHeader        `json:"images" form:"images"`
}

func (ui *UpdateService) GetServiceVersionLocalId() (uint32, error) {
	serviceVersionUID, err := common.FromBase58(ui.ServiceVersion.ID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return serviceVersionUID.GetLocalID(), nil
}

func (ui *UpdateService) GetServiceLocalId() (uint32, error) {
	serviceUID, err := common.FromBase58(ui.ServiceID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return serviceUID.GetLocalID(), nil
}

func (ui *UpdateServiceVersion) GetIntroVideoLocalId(ctx context.Context) (uint32, error) {
	introVideoUID, err := common.FromBase58(ui.IntroVideoID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid intro video ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid intro video ID"))
	}

	return introVideoUID.GetLocalID(), nil
}

func (ui *UpdateServiceVersion) GetCateogryLocalId(ctx context.Context) (uint32, error) {
	categoryUID, err := common.FromBase58(ui.CategoryID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid category ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid category ID"))
	}

	return categoryUID.GetLocalID(), nil
}

func (ui *UpdateServiceVersion) GetSubCategoryLocalId(ctx context.Context) (uint32, error) {
	subCategoryUID, err := common.FromBase58(ui.SubCategoryID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid sub category ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid sub category ID"))
	}

	return subCategoryUID.GetLocalID(), nil
}
