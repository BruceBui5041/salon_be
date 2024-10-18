package coursemodel

import (
	"mime/multipart"
	"video_server/common"
	"video_server/utils/customtypes"
)

type UpdateCourse struct {
	common.SQLModel `json:",inline"`
	Title           string                        `json:"title" form:"title"`
	Description     string                        `json:"description" form:"description"`
	CategoryID      string                        `json:"category_id" form:"category_id"`
	Thumbnail       *multipart.FileHeader         `json:"thumbnail" form:"thumbnail"`
	UploadedBy      string                        `json:"uploaded_by"`
	Price           customtypes.DecimalString     `json:"price" form:"price"`
	DiscountedPrice customtypes.NullDecimalString `json:"discounted_price" form:"discounted_price"`
	DifficultyLevel string                        `json:"difficulty_level" form:"difficulty_level"`
	Overview        string                        `json:"overview" form:"overview"`
	IntroVideoId    string                        `json:"intro_video_id" form:"intro_video_id"`
}

func (uc *UpdateCourse) Mask(isAdmin bool) {
	uc.GenUID(common.DbTypeCourse)
}
