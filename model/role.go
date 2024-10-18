package models

import (
	"video_server/common"
	"video_server/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const RoleEntityName = "Role"

func init() {
	modelhelper.RegisterModel(Role{})
}

type Role struct {
	common.SQLModel `json:",inline"`
	Name            string            `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Code            string            `json:"code" gorm:"uniqueIndex;not null;size:50"`
	Description     string            `json:"description"`
	Users           []*User           `json:"users,omitempty" gorm:"many2many:user_role;joinForeignKey:RoleID;joinReferences:UserID"`
	Permissions     []*Permission     `json:"permissions,omitempty" gorm:"many2many:role_permission;joinForeignKey:RoleID;joinReferences:PermissionID"`
	RolePermission  []*RolePermission `json:"role_permission,omitempty" gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return "role"
}

func (u *Role) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

func (r *Role) AfterFind(tx *gorm.DB) (err error) {
	r.Mask(false)
	return
}
