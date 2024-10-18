package models

import (
	"time"
	"video_server/common"
)

const EnrollmentEntityName = "Enrollment"

type Enrollment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32    `json:"user_id" gorm:"index"`
	CourseID        uint32    `json:"course_id" gorm:"index"`
	PaymentID       *uint32   `json:"payment_id,omitempty" gorm:"index"`
	EnrolledAt      time.Time `json:"enrolled_at" gorm:"autoCreateTime"`
	User            User      `json:"user" gorm:"constraint:OnDelete:CASCADE;"`
	Course          *Course   `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE;"`
	Payment         *Payment  `json:"payment,omitempty" gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL;"`
}

func (Enrollment) TableName() string {
	return "enrollment"
}

func (enr *Enrollment) Mask(isAdmin bool) {
	enr.GenUID(common.DBTypeEnrollment)
}
