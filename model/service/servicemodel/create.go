package servicemodel

import "salon_be/common"

type CreateService struct {
	common.SQLModel `json:",inline"`
	CreatorID       uint32                `json:"creator_id"`
	Slug            string                `json:"slug"`
	ServiceVersion  *CreateServiceVersion `json:"service_version,omitempty"`
}

type CreateServiceVersion struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	CategoryID      uint32   `json:"category_id"`
	IntroVideoID    *uint32  `json:"intro_video_id,omitempty"`
	Thumbnail       string   `json:"thumbnail"`
	Price           float64  `json:"price"`
	DiscountedPrice *float64 `json:"discounted_price,omitempty"`
}

func (cs *CreateService) Mask(isAdmin bool) {
	cs.GenUID(common.DBTypeService)
}
