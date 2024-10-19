package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const CategoryEntityName = "Category"

func init() {
	modelhelper.RegisterModel(Category{})
}

type Category struct {
	common.SQLModel `json:",inline"`
	Name            string            `json:"name" gorm:"not null;size:100"`
	Code            string            `json:"code" gorm:"not null;size:100"`
	Description     string            `json:"description"`
	Services        []*ServiceVersion `json:"services,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCategory)
}

func (c *Category) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	return
}
