package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init() {
	modelhelper.RegisterModel(OTP{})
}

const OTPEntityName = "OTP"

type OTP struct {
	common.SQLModel `json:",inline"`
	UUID            string     `json:"uuid" gorm:"column:uuid;type:varchar(50);uniqueIndex"`
	ESMSID          string     `json:"esmsid" gorm:"column:esmsid;type:varchar(50);not null"`
	UserID          uint32     `json:"-" gorm:"column:user_id;not null;index"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id;constraint:OnDelete:SET NULL;"`
	OTP             string     `json:"otp" gorm:"column:otp;type:varchar(6);uniqueIndex"`
	TTL             uint16     `json:"ttl" gorm:"column:ttl;type:smallint;not null"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"column:expires_at;type:timestamp"`
	PassedAt        *time.Time `json:"passed_at" gorm:"column:passed_at;type:timestamp"`
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

func (otp *OTP) BeforeCreate(tx *gorm.DB) error {
	otp.TTL = uint16(viper.GetInt("OTP_TTL"))
	return nil
}

func (otp *OTP) IsPassed(inputOTP string) bool {
	utcTime := time.Now().UTC()
	if inputOTP == otp.OTP && !utcTime.After(otp.ExpiresAt) {
		otp.PassedAt = &utcTime
		return true
	}
	return false
}
