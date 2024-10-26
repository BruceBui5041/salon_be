package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/storagehandler"

	"gorm.io/gorm"
)

const ImageEntityName = "Image"

func init() {
	modelhelper.RegisterModel(Comment{})
}

type Image struct {
	common.SQLModel `json:",inline"`
	UserID          uint32            `json:"user_id" gorm:"column:user_id;index"`
	User            *User             `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id;constraint:OnDelete:SET NULL;"`
	ServiceID       uint32            `json:"service_id" gorm:"column:service_id;index"`
	Service         *Service          `json:"service,omitempty" gorm:"foreignKey:ServiceID;references:Id;constraint:OnDelete:SET NULL;"`
	URL             string            `json:"url" gorm:"column:url;type:text"`
	ServiceVersions []*ServiceVersion `json:"service_versions,omitempty" gorm:"many2many:service_version_images;foreignKey:Id;joinForeignKey:ImageID;References:Id;joinReferences:ServiceVersionID;constraint:OnDelete:CASCADE;"`
}

func (Image) TableName() string {
	return "image"
}

func (i *Image) Mask(isAdmin bool) {
	i.GenUID(common.DBTypeComment)
}

func (i *Image) AfterFind(tx *gorm.DB) (err error) {
	i.Mask(false)
	i.URL = storagehandler.AddPublicCloudFrontDomain(i.URL)
	return
}
