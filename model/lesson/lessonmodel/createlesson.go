package lessonmodel

import (
	"video_server/common"
)

type CreateLesson struct {
	common.SQLModel `json:",inline"`
	CourseID        string `json:"course_id"`
	LectureID       string `json:"lecture_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Duration        int    `json:"duration"`
	Order           int    `json:"order"`
	Type            string `json:"type"`
	AllowPreview    bool   `json:"allow_preview"`
}

func (cl *CreateLesson) Mask(isAdmin bool) {
	cl.GenUID(common.DBTypeLesson)
}
