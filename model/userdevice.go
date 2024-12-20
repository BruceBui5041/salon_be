package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const UserDeviceEntityName = "UserDevice"

func init() {
	modelhelper.RegisterModel(UserDevice{})
}

type UserDevice struct {
	common.SQLModel `json:",inline"`
	FCMToken        string `json:"fcm_token" gorm:"column:fcm_token;uniqueIndex;type:varchar(250)"` // firebase cloud messaging token
	Platform        string `json:"platform" gorm:"column:platform;"`
	UserID          uint32 `json:"-" gorm:"column:user_id;not null;index"`
	User            *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (UserDevice) TableName() string {
	return "user_device"
}

func (u *UserDevice) Mask(isAdmin bool) {
	u.GenUID(common.DBTypeUserDevice)
}

func (u *UserDevice) AfterFind(tx *gorm.DB) (err error) {
	u.Mask(false)
	return
}
