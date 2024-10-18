package models

import (
	"fmt"
	"salon_be/common"

	"gorm.io/gorm"
)

const (
	LectureEntityName = "Lecture"
)

type Lecture struct {
	common.SQLModel `json:",inline"`
	CourseID        uint32    `json:"course_id" gorm:"index"`
	Title           string    `json:"title" gorm:"not null;size:255"`
	Description     string    `json:"description" gorm:"type:text"`
	Order           int       `json:"order"`
	Duration        int       `json:"duration"`
	Course          Course    `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL;"`
	Lessons         []*Lesson `json:"lessons" gorm:"foreignKey:LectureID;constraint:OnDelete:CASCADE;"`
}

func (Lecture) TableName() string {
	return "lecture"
}

func (l *Lecture) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLecture)
}

func (l *Lecture) BeforeCreate(tx *gorm.DB) (err error) {
	var maxOrder int
	sequence := 10
	if err := tx.Model(&Lecture{}).
		Where("course_id = ?", l.CourseID).
		Select("COALESCE(MAX(`order`), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	l.Order = maxOrder + sequence // Increment by 10 instead of 1

	// Use a transaction to prevent race conditions
	return tx.Transaction(func(tx *gorm.DB) error {
		// Check if the order is already taken (in case of concurrent inserts)
		var count int64
		if err := tx.Model(&Lecture{}).
			Where("course_id = ? AND `order` = ?", l.CourseID, l.Order).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			// If the order is taken, find the next available order
			return tx.Model(&Lecture{}).Where("course_id = ?", l.CourseID).
				Select(fmt.Sprintf("COALESCE(MAX(`order`), 0) + %d", sequence)).
				Scan(&l.Order).Error
		}

		return nil
	})
}

func (l *Lecture) AfterFind(tx *gorm.DB) (err error) {
	l.Mask(false)
	return
}
