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
	common.SQLModel  `json:",inline"`
	Name             string           `json:"name" gorm:"column:name;not null;size:255"`
	Description      string           `json:"description" gorm:"column:description;type:text"`
	Versions         []ServiceVersion `json:"versions,omitempty" gorm:"foreignKey:ServiceID"`
	ServiceVersion   *ServiceVersion  `json:"service_version,omitempty" gorm:"foreignKey:ActiveVersionID"`
	ServiceVersionID *uint32          `json:"service_version_id,omitempty" gorm:"column:service_version_id"`
}

func (Service) TableName() string {
	return "service"
}

func (s *Service) Mask(isAdmin bool) {
	s.GenUID(common.DBTypeService)
}

func (s *Service) AfterFind(tx *gorm.DB) (err error) {
	s.Mask(false)
	return
}
