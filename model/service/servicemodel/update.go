package servicemodel

import (
	"salon_be/common"
	"salon_be/utils/customtypes"
)

type UpdateService struct {
	ServiceID        string                `json:"service_id" form:"service_id"`
	ServiceVersionID string                `json:"service_version_id" form:"service_version_id"`
	ServiceVersion   *UpdateServiceVersion `json:"service_version,omitempty" form:"service_version,omitempty"`
}

func (ui *UpdateService) GetServiceVersionLocalId() (uint32, error) {
	serviceVersionUID, err := common.FromBase58(ui.ServiceVersionID)
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

type UpdateServiceVersion struct {
	Title           string                         `json:"title" form:"title"`
	Description     string                         `json:"description" form:"description"`
	CategoryID      string                         `json:"category_id" form:"category_id"`
	SubCategoryID   string                         `json:"sub_category_id" form:"sub_category_id"`
	IntroVideoID    string                         `json:"intro_video_id,omitempty" form:"intro_video_id,omitempty"`
	Thumbnail       string                         `json:"thumbnail" form:"thumbnail"`
	Price           customtypes.DecimalString      `json:"price" form:"price"`
	DiscountedPrice *customtypes.NullDecimalString `json:"discounted_price,omitempty" form:"discounted_price,omitempty"`
	Duration        uint32                         `json:"duration" form:"duration"`
}

func (ui *UpdateServiceVersion) GetCateogryLocalId() (uint32, error) {
	categoryUID, err := common.FromBase58(ui.CategoryID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return categoryUID.GetLocalID(), nil
}

func (ui *UpdateServiceVersion) GetSubCategoryLocalId() (uint32, error) {
	subCategoryUID, err := common.FromBase58(ui.SubCategoryID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return subCategoryUID.GetLocalID(), nil
}
