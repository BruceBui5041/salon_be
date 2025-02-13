package models

import (
	"salon_be/common"
	"salon_be/storagehandler"

	"gorm.io/gorm"
)

const CertificateEntityName = "Certificate"

type Certificate struct {
	common.SQLModel `json:",inline"`
	URL             string `json:"url" gorm:"column:url;size:255"`
	Type            string `json:"type" gorm:"column:type;size:50"`
	OwnerID         uint32 `json:"owner_id" gorm:"column:owner_id"`
	CreatorID       uint32 `json:"creator_id" gorm:"column:creator_id"`

	// Relations
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`
	Creator *User `json:"creator" gorm:"foreignKey:CreatorID"`
	// Images  []Image `json:"images" gorm:"many2many:certificate_images;"`
}

func (Certificate) TableName() string {
	return "certificates"
}

func (c *Certificate) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCertificate)
}

func (c *Certificate) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	c.URL = storagehandler.AddPublicCloudFrontDomain(c.URL)
	return
}
