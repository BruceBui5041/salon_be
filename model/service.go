package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const ServiceEntityName = "Service"

func init() {
	modelhelper.RegisterModel(Service{})
}

type Service struct {
	common.SQLModel `json:",inline"`
}

func (Service) TableName() string {
	return "service"
}

func (c *Service) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeService)
}

func (c *Service) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	return
}
