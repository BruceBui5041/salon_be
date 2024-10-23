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
	common.SQLModel  `json:",inline"`
	UserID           uint32          `json:"user_id" gorm:"column:user_id;index"`
	User             *User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
	ServiceID        uint32          `json:"service_id" gorm:"column:service_id;index"`
	Service          *Service        `json:"service,omitempty" gorm:"foreignKey:ServiceID;constraint:OnDelete:SET NULL;"`
	ServiceVersionID uint32          `json:"service_version_id" gorm:"column:service_version_id;index"`
	ServiceVersion   *ServiceVersion `json:"service_version,omitempty" gorm:"foreignKey:ServiceVersionID;constraint:OnDelete:SET NULL;"`
	URL              string          `json:"url" gorm:"column:url;type:text"`
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
