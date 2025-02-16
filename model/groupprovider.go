package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const GroupProviderEntityName = "GroupProvider"

func init() {
	modelhelper.RegisterModel(GroupProvider{})
}

type GroupProvider struct {
	common.SQLModel `json:",inline"`
	Name            string     `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Code            string     `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex;not null"`
	Description     string     `json:"description" gorm:"column:description;type:text"`
	OwnerID         uint32     `json:"-" gorm:"column:owner_id;index"`
	Owner           *User      `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:Id;constraint:OnDelete:SET NULL;"`
	CreatorID       uint32     `json:"-" gorm:"column:creator_id;index"`
	Creator         *User      `json:"creator,omitempty" gorm:"foreignKey:CreatorID;references:Id;constraint:OnDelete:SET NULL;"`
	Admins          []*User    `json:"admins,omitempty" gorm:"many2many:m2m_group_provider_admins;foreignKey:Id;joinForeignKey:GroupProviderID;References:Id;joinReferences:UserID;constraint:OnDelete:CASCADE;"`
	Providers       []*User    `json:"providers,omitempty" gorm:"many2many:m2m_group_provider_users;foreignKey:Id;joinForeignKey:GroupProviderID;References:Id;joinReferences:UserID;constraint:OnDelete:CASCADE;"`
	Services        []*Service `json:"services,omitempty" gorm:"foreignKey:GroupProviderID"`
	Images          []*Image   `json:"images,omitempty" gorm:"many2many:m2m_group_provider_images;foreignKey:Id;joinForeignKey:GroupProviderID;References:Id;joinReferences:ImageID;constraint:OnDelete:CASCADE;"`
}

func (GroupProvider) TableName() string {
	return "group_provider"
}

func (g *GroupProvider) Mask(isAdmin bool) {
	g.GenUID(common.DBTypeGroupProvider)
}

func (g *GroupProvider) AfterFind(tx *gorm.DB) (err error) {
	g.Mask(false)
	return
}
