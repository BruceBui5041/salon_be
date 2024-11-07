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

type UpdateVersionImage struct {
	ImageID string  `json:"image_id" form:"image_id"`
	Order   *uint32 `json:"order" form:"order"`
}

func (vi *UpdateVersionImage) GetLocalID(ctx context.Context) (uint32, error) {
	imageUID, err := common.FromBase58(vi.ImageID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid image ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid image ID"))
	}

	return imageUID.GetLocalID(), nil
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
	MainImageID     *string                        `json:"main_image_id,omitempty" form:"main_image_id,omitempty"`
	VersionImages   *[]UpdateVersionImage          `json:"version_images" form:"version_images"`
}

func (us *UpdateServiceVersion) GetMainImageLocalId(ctx context.Context) (uint32, error) {
	mainImageUID, err := common.FromBase58(*us.MainImageID)
	if err != nil {
		logger.AppLogger.Error(ctx, "invalid main image ID", zap.Error(err))
		return 0, common.ErrInvalidRequest(errors.New("invalid main image ID"))
	}

	return mainImageUID.GetLocalID(), nil
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
