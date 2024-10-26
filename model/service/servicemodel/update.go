package servicemodel

import "salon_be/utils/customtypes"

type UpdateService struct {
	ServiceID        uint32                `json:"service_id" form:"service_id"`
	ServiceVersionID uint32                `json:"service_version_id" form:"service_version_id"`
	ServiceVersion   *UpdateServiceVersion `json:"service_version,omitempty" form:"service_version,omitempty"`
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
