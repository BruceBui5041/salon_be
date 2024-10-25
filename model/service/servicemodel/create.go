package servicemodel

import (
	"salon_be/common"
	"salon_be/utils/customtypes"
)

type CreateService struct {
	common.SQLModel `json:",inline"`
	CreatorID       uint32                `json:"creator_id"`
	Slug            string                `json:"slug"`
	ServiceVersion  *CreateServiceVersion `json:"service_version,omitempty"`
}

type CreateServiceVersion struct {
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	CategoryID      string                         `json:"category_id"`
	SubCategoryID   string                         `json:"sub_category_id"`
	IntroVideoID    *string                        `json:"intro_video_id,omitempty"`
	Thumbnail       string                         `json:"thumbnail"`
	Price           customtypes.DecimalString      `json:"price"`
	DiscountedPrice *customtypes.NullDecimalString `json:"discounted_price,omitempty"`
}

func (cs *CreateService) Mask(isAdmin bool) {
	cs.GenUID(common.DBTypeService)
}
