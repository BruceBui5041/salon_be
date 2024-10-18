package videomodel

import (
	"salon_be/common"
)

type CreateVideo struct {
	common.SQLModel `json:",inline"`
	ServiceID       string `json:"service_id" form:"service_id"`
	Title           string `json:"title" form:"title"`
	Description     string `json:"description" form:"description"`
	VideoURL        string `json:"video_url" form:"video_url"`
	ThumbnailURL    string `json:"thumbnail_url" form:"thumbnail_url"`
	Duration        int    `json:"duration" form:"duration"`
	Order           int    `json:"order" form:"order"`
}

func (cv *CreateVideo) Mask(isAdmin bool) {
	cv.GenUID(common.DbTypeVideo)
}
