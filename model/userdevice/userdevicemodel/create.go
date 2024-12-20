package userdevicemodel

import (
	"salon_be/common"
)

type CreateUserDevice struct {
	common.SQLModel `json:",inline"`
	FCMToken        string `json:"fcm_token" form:"fcm_token"`
	Platform        string `json:"platform" form:"platform"`
	UserID          uint32 `json:"-"`
}

func (CreateUserDevice) TableName() string {
	return "user_device"
}
