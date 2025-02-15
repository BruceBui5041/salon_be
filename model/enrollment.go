package models

import (
	"salon_be/common"
	"time"
)

const EnrollmentEntityName = "Enrollment"

type Enrollment struct {
	common.SQLModel  `json:",inline"`
	UserID           uint32          `json:"user_id" gorm:"index"`
	ServiceVersionID uint32          `json:"service_version_id" gorm:"index"`
	PaymentID        *uint32         `json:"payment_id,omitempty" gorm:"index"`
	EnrolledAt       time.Time       `json:"enrolled_at" gorm:"autoCreateTime"`
	User             User            `json:"user" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	ServiceVersion   *ServiceVersion `json:"service,omitempty" gorm:"foreignKey:ServiceVersionID;constraint:OnDelete:CASCADE;"`
	Payment          *Payment        `json:"payment,omitempty" gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL;"`
}

func (Enrollment) TableName() string {
	return "enrollment"
}

func (enr *Enrollment) Mask(isAdmin bool) {
	enr.GenUID(common.DBTypeEnrollment)
}
