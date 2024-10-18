package lessonmodel

import (
	"video_server/common"
	models "video_server/model"
)

type LessonResponse struct {
	common.SQLModel `json:",inline"`
	CourseID        uint32 `json:"-"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Duration        int    `json:"duration"`
	Order           int    `json:"order"`
	Type            string `json:"type"`
	AllowPreview    bool   `json:"allow_preview"`
	// Videos          []videomodel.GetCourseVideoReponse `json:"videos"`
}

func (LessonResponse) TableName() string {
	return models.Lesson{}.TableName()
}

func (l *LessonResponse) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLesson)
}
