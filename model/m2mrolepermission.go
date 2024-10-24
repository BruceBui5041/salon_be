package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"time"

	"gorm.io/gorm"
)

const RolePermissionEntityName = "RolePermission"

func init() {
	modelhelper.RegisterModel(RolePermission{})
}

type RolePermission struct {
	RoleID           uint32     `gorm:"primaryKey;column:role_id" json:"-"`
	PermissionID     uint32     `gorm:"primaryKey;column:permission_id" json:"-"`
	FakeRoleID       string     `gorm:"primaryKey;column:role_id" json:"role_id"`
	FakePermissionID string     `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	CreatedAt        *time.Time `json:"created_at,omitempty" gorm:"column,created_at;"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" gorm:"column,updated_at;"`
	CreatePermission bool       `gorm:"column:create_permission;default:0" json:"create_permission"`
	ReadPermission   bool       `gorm:"column:read_permission;default:0" json:"read_permission"`
	WritePermission  bool       `gorm:"column:write_permission;default:0" json:"write_permission"`
	DeletePermission bool       `gorm:"column:delete_permission;default:0" json:"delete_permission"`

	Role       *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}

func (rp *RolePermission) AfterFind(tx *gorm.DB) (err error) {
	tempPermission := Permission{SQLModel: common.SQLModel{Id: rp.PermissionID}}
	tempPermission.Mask(false)
	rp.FakePermissionID = tempPermission.GetFakeId()

	tempRole := Role{SQLModel: common.SQLModel{Id: rp.RoleID}}
	tempRole.Mask(false)
	rp.FakeRoleID = tempRole.GetFakeId()

	return
}
