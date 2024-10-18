package models

import (
	"salon_be/common"
	"salon_be/storagehandler"

	"gorm.io/gorm"
)

const (
	VideoEntityName = "Video"
)

type Video struct {
	common.SQLModel `json:",inline"`
	ServiceID       uint32              `json:"service_id" gorm:"index"`
	LessonID        *uint32             `json:"lesson_id" gorm:"index"`
	Title           string              `json:"title" gorm:"not null;size:255"`
	Description     string              `json:"description"`
	VideoURL        string              `json:"video_url" gorm:"not null;size:255"`
	RawVideoURL     string              `json:"raw_video_url" gorm:"not null;size:255"`
	ThumbnailURL    string              `json:"thumbnail_url" gorm:"not null;size:255"`
	Duration        int                 `json:"duration"`
	Order           int                 `json:"order"`
	Service         Service             `json:"service" gorm:"constraint:OnDelete:CASCADE;"`
	Tags            []*Tag              `json:"tags,omitempty" gorm:"many2many:video_tags;"`
	Progress        []*Progress         `json:"progress,omitempty"`
	AllowPreview    bool                `json:"allow_preview" gorm:"default:false"`
	ProcessInfos    []*VideoProcessInfo `json:"process_infos,omitempty" gorm:"foreignKey:VideoID"`
}

func (Video) TableName() string {
	return "video"
}

func (v *Video) Mask(isAdmin bool) {
	v.GenUID(common.DbTypeVideo)
}

func (v *Video) AfterFind(tx *gorm.DB) (err error) {
	v.Mask(false)
	v.ThumbnailURL = storagehandler.AddPublicCloudFrontDomain(v.ThumbnailURL)
	return
}
