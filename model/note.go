package models

import (
	"video_server/common"

	"gorm.io/gorm"
)

const NoteEntityName = "Note"

type Note struct {
	common.SQLModel `json:",inline"`
	UserID          uint32  `json:"user_id" gorm:"column:user_id;index"`
	LessonID        uint32  `json:"lesson_id" gorm:"column:lesson_id;index"`
	User            *User   `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
	Lesson          *Lesson `json:"lesson,omitempty" gorm:"foreignKey:LessonID;constraint:OnDelete:SET NULL;"`
	Content         string  `json:"content" gorm:"column:content;type:text"`
	TimeMarked      string  `json:"time_marked" gorm:"column:content;type:varchar(20)"`
}

func (Note) TableName() string {
	return "note"
}

func (n *Note) Mask(isAdmin bool) {
	n.GenUID(common.DBTypeNote)
}

func (n *Note) AfterFind(tx *gorm.DB) (err error) {
	n.Mask(false)
	return
}
