package lecturemodel

import (
	"video_server/common"
)

type CreateLecture struct {
	common.SQLModel `json:",inline"`
	CourseID        string `json:"course_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Order           int    `json:"order"`
}

func (cl *CreateLecture) Mask(isAdmin bool) {
	cl.GenUID(common.DBTypeLecture)
}
