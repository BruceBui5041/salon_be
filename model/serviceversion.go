package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/storagehandler"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const ServiceVersionEntityName = "Service Version"

func init() {
	modelhelper.RegisterModel(ServiceVersion{})
}

type ReviewInfo struct {
	Stars uint8 `json:"stars"`
	Count uint  `json:"count"`
}

type ServiceVersion struct {
	common.SQLModel `json:",inline"`
	ServiceID       uint32               `json:"-" gorm:"column:service_id;not null;index"`
	Service         *Service             `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	IntroVideoID    *uint32              `json:"-" gorm:"column:intro_video_id;index"`
	IntroVideo      *Video               `json:"intro_video,omitempty" gorm:"foreignKey:IntroVideoID"`
	CategoryID      uint32               `json:"-" gorm:"column:category_id;index"`
	Category        *Category            `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL;"`
	SubCategoryID   uint32               `json:"-" gorm:"column:sub_category_id;index"`
	SubCategory     *Category            `json:"sub_category,omitempty" gorm:"foreignKey:SubCategoryID;constraint:OnDelete:SET NULL;"`
	Enrollments     []*Enrollment        `json:"enrollments,omitempty" gorm:"foreignKey:ServiceVersionID"`
	Images          []Image              `json:"images,omitempty" gorm:"foreignKey:ServiceVersionID"`
	Title           string               `json:"title" gorm:"column:title;not null;size:255"`
	Description     string               `json:"description" gorm:"column:description;type:text"`
	ServiceMen      []User               `json:"service_men" gorm:"many2many:user_service;"`
	Thumbnail       string               `json:"thumbnail" gorm:"column:thumbnail;type:text"`
	Price           decimal.Decimal      `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	DiscountedPrice *decimal.NullDecimal `json:"discounted_price" gorm:"column:discounted_price;type:decimal(10,2)"`
}

func (ServiceVersion) TableName() string {
	return "service_version"
}

func (sv *ServiceVersion) Mask(isAdmin bool) {
	sv.GenUID(common.DbTypeServiceVersion)
}

func (sv *ServiceVersion) AfterFind(tx *gorm.DB) (err error) {
	sv.Mask(false)
	if sv.Thumbnail != "" {
		sv.Thumbnail = storagehandler.AddPublicCloudFrontDomain(sv.Thumbnail)
	}
	return
}
