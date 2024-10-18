package models

import (
	"video_server/common"
	"video_server/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const CategoryEntityName = "Category"

func init() {
	modelhelper.RegisterModel(Category{})
}

type Category struct {
	common.SQLModel `json:",inline"`
	Name            string    `json:"name" gorm:"not null;size:100"`
	Description     string    `json:"description"`
	Courses         []*Course `json:"courses,omitempty" gorm:"foreignKey:CategoryID"`
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
