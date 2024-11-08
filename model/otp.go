package models

import (
	"salon_be/common"
	"time"

	"gorm.io/gorm"
)

const OTPEntityName = "OTP"

type OTP struct {
	common.SQLModel `json:",inline"`
	UserID          uint32    `json:"-" gorm:"column:user_id;not null;index"`
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id;constraint:OnDelete:SET NULL;"`
	OTP             string    `json:"otp" gorm:"column:otp;type:varchar(6);uniqueIndex"`
	TTL             uint16    `json:"ttl" gorm:"column:ttl;type:smallint;not null"`
	ExpiresAt       time.Time `json:"expires_at" gorm:"column:expires_at;type:timestamp"`
	PassedAt        time.Time `json:"passed_at" gorm:"column:passed_at;type:timestamp"`
}

func (OTP) TableName() string {
	return "otp"
}

func (otp *OTP) Mask(isAdmin bool) {
	otp.GenUID(common.DBTypePayment)
}

func (otp *OTP) AfterFind(tx *gorm.DB) (err error) {
	otp.Mask(false)
	return
}

func (otp *OTP) IsPassed(inputOTP string, now time.Time) bool {
	if inputOTP == otp.OTP && !now.After(otp.ExpiresAt) {
		otp.PassedAt = now
		return true
	}
	return false
}
