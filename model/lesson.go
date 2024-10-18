package models

import (
	"fmt"
	"video_server/common"

	"gorm.io/gorm"
)

const (
	LessonEntityName = "Lesson"
)

type Lesson struct {
	common.SQLModel `json:",inline"`
	Type            string  `json:"type" gorm:"column:type;type:ENUM('video','quiz','assignment');default:video"`
	CourseID        uint32  `json:"course_id" gorm:"index"`
	LectureID       uint32  `json:"lecture_id" gorm:"index"`
	Title           string  `json:"title" gorm:"not null;size:255"`
	Description     string  `json:"description" gorm:"type:text"`
	Duration        int     `json:"duration"`
	Order           int     `json:"order"`
	Videos          []Video `json:"videos" gorm:"foreignKey:LessonID"`
	Lecture         Lecture `json:"lecture" gorm:"foreignKey:LectureID;constraint:OnDelete:SET NULL;"`
	Course          Course  `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL;"`
	AllowPreview    bool    `json:"allow_preview"`
}

func (Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLesson)
}

func (l *Lesson) AfterFind(tx *gorm.DB) (err error) {
	l.Mask(false)
	return
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) (err error) {
	var maxOrder int
	sequence := 10

	// Find the maximum order for lessons within the same lecture
	if err := tx.Model(&Lesson{}).
		Where("lecture_id = ?", l.LectureID).
		Select("COALESCE(MAX(`order`), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	l.Order = maxOrder + sequence // Increment by 10 instead of 1

	// Use a transaction to prevent race conditions
	return tx.Transaction(func(tx *gorm.DB) error {
		// Check if the order is already taken (in case of concurrent inserts)
		var count int64
		if err := tx.Model(&Lesson{}).
			Where("lecture_id = ? AND `order` = ?", l.LectureID, l.Order).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			// If the order is taken, find the next available order
			return tx.Model(&Lesson{}).
				Where("lecture_id = ?", l.LectureID).
				Select(fmt.Sprintf("COALESCE(MAX(`order`), 0) + %d", sequence)).
				Scan(&l.Order).Error
		}

		return nil
	})
}
