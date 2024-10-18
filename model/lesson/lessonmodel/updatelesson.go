package lessonmodel

import (
	"video_server/common"
	models "video_server/model"
)

type UpdateLesson struct {
	common.SQLModel `json:",inline"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Duration        int     `json:"duration"`
	Order           int     `json:"order"`
	VideoID         *string `json:"video_id" gorm:"-"`
	Type            string  `json:"type"`
	AllowPreview    bool    `json:"allow_preview"`
}

func (UpdateLesson) TableName() string {
	return models.Lesson{}.TableName()
}

func (ul *UpdateLesson) Mask(isAdmin bool) {
	// No need to mask anything for update
}
