// File: model/lecture/lecturemodel/updatelecture.go

package lecturemodel

import (
	"video_server/common"
	models "video_server/model"
)

type UpdateLecture struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Order           int    `json:"order"`
}

func (UpdateLecture) TableName() string {
	return models.Lecture{}.TableName()
}

func (ul *UpdateLecture) Mask(isAdmin bool) {
	// No need to mask anything for update
}
