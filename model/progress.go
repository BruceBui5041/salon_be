package models

import (
	"time"
	"video_server/common"

	"gorm.io/gorm"
)

const ProgressEntityName = "Progress"

type Progress struct {
	common.SQLModel `json:",inline"`
	UserID          uint32    `json:"user_id" gorm:"column:user_id;index"`
	VideoID         uint32    `json:"video_id" gorm:"column:video_id;index"`
	WatchedSeconds  uint32    `json:"watched_seconds" gorm:"column:watched_seconds;default:0"`
	Completed       bool      `json:"completed" gorm:"column:completed;default:false"`
	LastWatched     time.Time `json:"last_watched" gorm:"column:last_watched;autoUpdateTime"`
	User            User      `json:"user" gorm:"constraint:OnDelete:CASCADE;"`
	Video           Video     `json:"video" gorm:"constraint:OnDelete:CASCADE;"`
}

func (Progress) TableName() string {
	return "progress"
}

func (p *Progress) Mask(isAdmin bool) {
	p.GenUID(common.DBTypePayment)
}

func (p *Progress) AfterFind(tx *gorm.DB) (err error) {
	p.Mask(false)
	return
}
