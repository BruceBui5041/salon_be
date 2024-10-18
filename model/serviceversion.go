package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	ServiceID       uint32          `json:"service_id" gorm:"column:service_id;not null;index"`
	Service         *Service        `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	Title           string          `json:"title" gorm:"column:title;not null;size:255"`
	Description     string          `json:"description" gorm:"column:description;type:text"`
	CreatorID       uint32          `json:"creator_id" gorm:"column:creator_id;index"`
	CategoryID      uint32          `json:"category_id" gorm:"column:category_id;index"`
	Creator         *User           `json:"creator,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	Category        *Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL;"`
	Instructors     []User          `json:"instructors" gorm:"many2many:user_service;"`
	IntroVideo      *Video          `json:"intro_video,omitempty" gorm:"foreignKey:IntroVideoID"`
	Enrollments     []Enrollment    `json:"enrollments,omitempty" gorm:"foreignKey:ServiceVersionID"`
	Slug            string          `json:"slug" gorm:"column:slug;not null;size:255"`
	Thumbnail       string          `json:"thumbnail" gorm:"column:thumbnail;type:text"`
	Price           decimal.Decimal `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	Comments        []*Comment      `json:"comments,omitempty" gorm:"foreignKey:ServiceID"`
	RatingCount     uint            `json:"rating_count" gorm:"column:rating_count;type:unsigned int"`
	ReviewInfo      ReviewInfos     `json:"review_info" gorm:"column:review_info;type:json"`
	AverageRating   decimal.Decimal `json:"avg_rating" gorm:"column:avg_rating;type:decimal(3,1)"`
}

func (ServiceVersion) TableName() string {
	return "service_version"
}

func (sv *ServiceVersion) Mask(isAdmin bool) {
	sv.GenUID(common.DbTypeServiceVersion)
}

func (sv *ServiceVersion) AfterFind(tx *gorm.DB) (err error) {
	sv.Mask(false)
	sv.Thumbnail = storagehandler.AddPublicCloudFrontDomain(sv.Thumbnail)
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
