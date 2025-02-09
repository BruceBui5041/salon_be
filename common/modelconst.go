package common

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	DbTypeVideo              = 1
	DbTypeServiceVersion     = 2
	DbTypeTag                = 3
	DbTypeUser               = 4
	DBTypeCategory           = 5
	DBTypeService            = 6
	DBTypeVideoProcessInfo   = 8
	DBTypeEnrollment         = 9
	DBTypePayment            = 10
	DBTypeUserProfile        = 11
	DBTypeComment            = 12
	DBTypeRate               = 13
	DBTypeNote               = 14
	DBTypePermission         = 15
	DBTypeProgress           = 16
	DBTypeBooking            = 17
	DBTypeCoupon             = 18
	DBTypeNotification       = 19
	DBTypeNotificationDetail = 20
	DBTypeUserDevice         = 21
	DBTypeCommission         = 22
	DBTypeLocation           = 23
	DBTypeBookingLocation    = 24
	DBTypeDistanceTracking   = 25
)

const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusSuspended = "suspended"
)

const CurrentUser = "user"

type PreloadInfo struct {
	Name     string
	Function func(*gorm.DB) *gorm.DB
}

type Requester interface {
	GetUserId() uint32
	GetEmail() string
	GetRoles(ctx context.Context) []byte
	GetFakeId() string
	Mask(isAdmin bool)
	IsAdmin() bool
	IsUser() bool
	IsSuperAdmin() bool
	IsProvider() bool
}

type SQLModel struct {
	Id        uint32     `json:"-" gorm:"column,id;" form:"id"`
	FakeId    *UID       `json:"id" gorm:"-"`
	Status    string     `json:"status" gorm:"column:status;type:ENUM('active','inactive','suspended');default:active" form:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column,created_at;" form:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column,updated_at;" form:"updated_at"`
}

func (model *SQLModel) GenUID(dbType int) {
	uid := NewUID(uint32(model.Id), dbType, 1)
	model.FakeId = &uid
}

func (model *SQLModel) GetFakeId() string {
	return model.FakeId.String()
}
