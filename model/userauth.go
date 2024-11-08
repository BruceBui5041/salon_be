package models

import "salon_be/common"

type UserAuth struct {
	common.SQLModel         `json:",inline"`
	UserID                  uint32 `gorm:"column:user_id;index"`
	AuthType                string `gorm:"column:auth_type;not null;size:20"`
	AuthProviderID          string `gorm:"column:auth_provider_id;size:255"`
	AuthProviderToken       string `gorm:"column:auth_provider_token"`
	User                    User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;"`
	PhoneNumberVerifyStatus string `gorm:"column:phone_number_verify_status;not null;type:ENUM('verified','unverified');default:unverified;size:20"`
}

func (UserAuth) TableName() string {
	return "user_auth"
}
