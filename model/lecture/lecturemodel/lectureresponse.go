package lecturemodel

import (
	"video_server/common"
	models "video_server/model"
	"video_server/model/lesson/lessonmodel"
)

const (
	LectureEntityName = "Lecture"
)

type LectureResponse struct {
	common.SQLModel `json:",inline"`
	CourseID        uint32                       `json:"-" gorm:"index"`
	Title           string                       `json:"title" gorm:"not null;size:255"`
	Description     string                       `json:"description" gorm:"type:text"`
	Order           int                          `json:"order"`
	Course          models.Course                `json:"-" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL;"`
	Lessons         []lessonmodel.LessonResponse `json:"lessons" gorm:"foreignKey:LectureID"`
	Duration        int                          `json:"duration"`
}

func (l *LectureResponse) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLecture)
}
