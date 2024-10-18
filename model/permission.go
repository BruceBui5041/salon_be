package models

import (
	"video_server/common"
	"video_server/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const PermissionEntityName = "Permission"

func init() {
	modelhelper.RegisterModel(Permission{})
}

type Permission struct {
	common.SQLModel `json:",inline"`
	Name            string            `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Code            string            `json:"code" gorm:"uniqueIndex;not null;size:50"`
	Description     string            `json:"description"`
	Roles           []*Role           `json:"roles,omitempty" gorm:"many2many:role_permission;joinForeignKey:PermissionID;joinReferences:RoleID"`
	RolePermission  []*RolePermission `json:"role_permission,omitempty" gorm:"foreignKey:PermissionID"`
}

func (Permission) TableName() string {
	return "permission"
}

func (p *Permission) Mask(isAdmin bool) {
	p.GenUID(common.DBTypePermission)
}

func (p *Permission) AfterFind(tx *gorm.DB) (err error) {
	p.Mask(false)
	return
}
