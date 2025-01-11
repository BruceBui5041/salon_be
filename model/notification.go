package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"time"

	"gorm.io/gorm"
)

const (
	NotificationEntityName = "Notification"

	// Notification states
	NotificationStatePending = "pending"
	NotificationStateSent    = "sent"
	NotificationStateError   = "error"
)

func init() {
	modelhelper.RegisterModel(Notification{})
}

// Metadata type for metadata
type Metadata map[string]interface{}

func (j Metadata) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *Metadata) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}
	return json.Unmarshal(bytes, &j)
}

// Notification model
type Notification struct {
	common.SQLModel `json:",inline"`
	Type            string                `json:"type" gorm:"column:type;type:varchar(50);index"`
	Scheduled       *time.Time            `json:"scheduled" gorm:"column:scheduled;type:datetime"`
	Metadata        Metadata              `json:"metadata" gorm:"column:metadata;type:json"`
	Details         []*NotificationDetail `json:"details" gorm:"foreignKey:NotificationID"`
	BookingID       uint32                `json:"-" gorm:"column:booking_id;index"`
	Booking         *Booking              `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.Status = common.StatusActive
	return nil
}

func (n *Notification) Mask(isAdmin bool) {
	n.GenUID(common.DBTypeNotification)
}

func (n *Notification) AfterFind(tx *gorm.DB) (err error) {
	n.Mask(false)
	return
}
