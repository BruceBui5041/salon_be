package servicemodel

import (
	"mime/multipart"
	"salon_be/common"
	"salon_be/utils/customtypes"
)

type CreateServiceRequest struct {
	JSON   string                  `form:"json"`
	Images []*multipart.FileHeader `form:"images"`
}

type CreateService struct {
	common.SQLModel `json:",inline"`
	CreatorID       uint32                `json:"creator_id" form:"creator_id"`
	Slug            string                `json:"slug" form:"slug"`
	ServiceVersion  *CreateServiceVersion `json:"service_version" form:"service_version"`
}

type CreateServiceVersion struct {
	Title           string                         `json:"title" form:"title"`
	Description     string                         `json:"description" form:"description"`
	CategoryID      string                         `json:"category_id" form:"category_id"`
	SubCategoryID   string                         `json:"sub_category_id" form:"sub_category_id"`
	IntroVideoID    *string                        `json:"intro_video_id,omitempty" form:"intro_video_id"`
	Thumbnail       string                         `json:"thumbnail" form:"thumbnail"`
	Price           customtypes.DecimalString      `json:"price" form:"price"`
	DiscountedPrice *customtypes.NullDecimalString `json:"discounted_price,omitempty" form:"discounted_price"`
	Duration        uint32                         `json:"duration" form:"duration"`
	Images          []*multipart.FileHeader        `json:"images" form:"images"`
}

func (cs *CreateService) Mask(isAdmin bool) {
	cs.GenUID(common.DBTypeService)
}
