package videomodel

import (
	"salon_be/common"
)

type GetCourseVideoReponse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title" `
	Slug            string `json:"slug" `
	Description     string `json:"description"`
	VideoURL        string `json:"video_url" `
	ThumbnailURL    string `json:"thumbnail_url" `
	Duration        int    `json:"duration"`
	Order           int    `json:"order"`
	AllowPreview    bool   `json:"allow_preview" `
	Overview        string `json:"overview"`
}

func (g *GetCourseVideoReponse) Mask(isAdmin bool) {
	g.GenUID(common.DbTypeVideo)
}
