package coursemodel

import (
	"mime/multipart"
	"video_server/common"
	"video_server/utils/customtypes"
)

// CreateUser represents the data needed to create a new user
type CreateCourse struct {
	common.SQLModel `json:",inline"`
	Title           string                        `json:"title" form:"title"`
	Description     string                        `json:"description" form:"description"`
	CategoryID      string                        `json:"category_id" form:"category_id"`
	CreatorID       uint32                        `json:"creator_id" form:"creator_id"`
	Slug            string                        `json:"slug" form:"slug"`
	Thumbnail       *multipart.FileHeader         `json:"thumbnail" form:"thumbnail"`
	Price           customtypes.DecimalString     `json:"price" form:"price"`
	DiscountedPrice customtypes.NullDecimalString `json:"discounted_price" form:"discounted_price"`
	DifficultyLevel string                        `json:"difficulty_level" form:"difficulty_level"`
}

func (cc *CreateCourse) Mask(isAdmin bool) {
	cc.GenUID(common.DbTypeCourse)
}
