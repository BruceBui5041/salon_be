package models

import (
	"video_server/common"
)

type LessonVideo struct {
	LessonID uint32 `gorm:"primaryKey"`
	VideoID  uint32 `gorm:"primaryKey"`
	common.SQLModel
}

func (LessonVideo) TableName() string {
	return "m2m_lesson_video"
}
