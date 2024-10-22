package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/storagehandler"

	"gorm.io/gorm"
)

const CategoryEntityName = "Category"

func init() {
	modelhelper.RegisterModel(Category{})
}

type Category struct {
	common.SQLModel `json:",inline"`
	Name            string            `json:"name" gorm:"not null;size:100"`
	Code            string            `json:"code" gorm:"not null;size:100;unique"`
	Image           string            `json:"image" gorm:"not null;size:255"`
	OriginImage     string            `json:"-" gorm:"-"`
	Description     string            `json:"description"`
	ParentID        *uint32           `json:"-" gorm:"column:parent_id;default:null"`
	Parent          *Category         `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	SubCategories   []*Category       `json:"sub_categories,omitempty" gorm:"foreignKey:ParentID"`
	ServiceVersions []*ServiceVersion `json:"service_versions,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCategory)
}

func (c *Category) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	c.Image = storagehandler.AddPublicCloudFrontDomain(c.Image)
	c.OriginImage = c.Image
	return
}
