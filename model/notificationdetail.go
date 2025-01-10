package models

import (
	"salon_be/common"
	"time"

	"gorm.io/gorm"
)

// NotificationDetail model (many-to-many relationship)
type NotificationDetail struct {
	common.SQLModel `json:",inline"`
	NotificationID  uint32     `json:"notification_id" gorm:"column:notification_id;index"`
	UserID          uint32     `json:"user_id" gorm:"column:user_id;index"`
	State           string     `json:"state" gorm:"column:state;type:ENUM('pending','sent','error');default:pending"`
	SentAt          *time.Time `json:"sent_at,omitempty" gorm:"column:sent_at"`
	Error           string     `json:"error,omitempty" gorm:"column:error;type:text"`
	ReadAt          *time.Time `json:"read_at,omitempty" gorm:"column:read_at"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	MessageID       string     `json:"message_id,omitempty" gorm:"column:message_id"`
}

func (NotificationDetail) TableName() string {
	return "notification_details"
}

func (nd *NotificationDetail) BeforeCreate(tx *gorm.DB) error {
	nd.Status = common.StatusActive
	nd.State = NotificationStatePending
	return nil
}

func (nd *NotificationDetail) Mask(isAdmin bool) {
	nd.GenUID(common.DBTypeNotificationDetail) // You'll need to add this constant to common/modelconst.go
}

func (nd *NotificationDetail) AfterFind(tx *gorm.DB) (err error) {
	nd.Mask(false)
	return
}
