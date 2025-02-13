package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/storagehandler"
	"time"

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
	Bookings        []*Booking           `json:"bookings,omitempty" gorm:"many2many:m2mbooking_service_version;foreignKey:Id;joinForeignKey:ServiceVersionID;References:Id;joinReferences:BookingID"`
	ServiceID       uint32               `json:"-" gorm:"column:service_id;not null;index"`
	Service         *Service             `json:"service,omitempty" gorm:"foreignKey:ServiceID;references:Id"`
	IntroVideoID    *uint32              `json:"-" gorm:"column:intro_video_id;index"`
	IntroVideo      *Video               `json:"intro_video,omitempty" gorm:"foreignKey:IntroVideoID;references:Id"`
	CategoryID      uint32               `json:"-" gorm:"column:category_id;index"`
	Category        *Category            `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:Id;constraint:OnDelete:SET NULL;"`
	SubCategoryID   uint32               `json:"-" gorm:"column:sub_category_id;index"`
	SubCategory     *Category            `json:"sub_category,omitempty" gorm:"foreignKey:SubCategoryID;references:Id;constraint:OnDelete:SET NULL;"`
	Enrollments     []*Enrollment        `json:"enrollments,omitempty" gorm:"foreignKey:ServiceVersionID;references:Id"`
	Images          []*Image             `json:"images,omitempty" gorm:"many2many:service_version_images;foreignKey:Id;joinForeignKey:ServiceVersionID;References:Id;joinReferences:ImageID;constraint:OnDelete:CASCADE;"`
	MainImageID     uint32               `json:"-" gorm:"column:main_image_id;index"`
	MainImage       *Image               `json:"main_image,omitempty" gorm:"foreignKey:MainImageID;references:Id;constraint:OnDelete:SET NULL"`
	Title           string               `json:"title" gorm:"column:title;not null;size:255"`
	Description     string               `json:"description" gorm:"column:description;type:text"`
	ServiceMen      []*User              `json:"service_men,omitempty" gorm:"many2many:user_service;foreignKey:Id;joinForeignKey:ServiceVersionID;References:Id;joinReferences:UserID"`
	Thumbnail       string               `json:"thumbnail" gorm:"column:thumbnail;type:text"`
	Price           decimal.Decimal      `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	DiscountedPrice *decimal.NullDecimal `json:"discounted_price" gorm:"column:discounted_price;type:decimal(10,2)"`
	Duration        uint32               `json:"duration" gorm:"column:duration;not null"`
	PublishedDate   *time.Time           `json:"published_date" gorm:"column:published_date;type:datetime"`
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
