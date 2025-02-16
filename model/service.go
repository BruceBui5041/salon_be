package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const ServiceEntityName = "Service"

func init() {
	modelhelper.RegisterModel(Service{})
}

type Service struct {
	common.SQLModel  `json:",inline"`
	CreatorID        uint32           `json:"-" gorm:"column:creator_id;index"`
	Creator          *User            `json:"creator,omitempty" gorm:"constraint:OnDelete:SET NULL;foreignKey:CreatorID"`
	OwnerID          uint32           `json:"-" gorm:"column:owner_id;index"`
	Owner            *User            `json:"owner,omitempty" gorm:"constraint:OnDelete:SET NULL;foreignKey:OwnerID"`
	Versions         []ServiceVersion `json:"versions,omitempty" gorm:"foreignKey:ServiceID"`
	ServiceVersion   *ServiceVersion  `json:"service_version,omitempty" gorm:"foreignKey:ServiceVersionID"`
	ServiceVersionID *uint32          `json:"-" gorm:"column:service_version_id"`
	GroupProviderID  *uint32          `json:"-" gorm:"column:group_provider_id"`
	GroupProvider    *GroupProvider   `json:"group_provider,omitempty" gorm:"constraint:OnDelete:SET NULL;foreignKey:GroupProviderID"`
	Comments         []*Comment       `json:"comments,omitempty" gorm:"foreignKey:ServiceID"`
	Slug             string           `json:"slug" gorm:"column:slug;not null;size:255"`
	RatingCount      uint32           `json:"rating_count" gorm:"column:rating_count"`
	ReviewInfo       ReviewInfos      `json:"review_info" gorm:"column:review_info;type:json"`
	AverageRating    decimal.Decimal  `json:"avg_rating" gorm:"column:avg_rating;type:decimal(3,1)"`
	Images           []Image          `json:"images,omitempty" gorm:"foreignKey:ServiceID"`
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

type ReviewInfos []ReviewInfo

func (r *ReviewInfos) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, r)
}

func (r ReviewInfos) Value() (driver.Value, error) {
	return json.Marshal(r)
}
