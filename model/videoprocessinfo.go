package models

import (
	"salon_be/common"
)

const (
	VideoProgressInfoEntityName = "VideoProcessInfo"
)

type VideoProcessInfo struct {
	common.SQLModel   `json:",inline"`
	VideoID           uint32 `json:"video_id" gorm:"not null;index"`
	ProcessResolution string `json:"process_resolution" gorm:"type:ENUM('360p','480p','720p','1080p')"`
	ProcessState      string `json:"process_state" gorm:"type:ENUM('pending','inqueue','processing','done','error');default:pending"`
	ProcessError      string `json:"process_error" gorm:"type:text"`
	ProcessRetry      uint16 `json:"process_retry" gorm:"type:smallint unsigned;default:0"`
	ProcessProgress   uint16 `json:"process_progress" gorm:"type:smallint unsigned;default:0"`
	Video             Video  `json:"video" gorm:"foreignKey:VideoID;constraint:OnDelete:CASCADE;"`
}

func (VideoProcessInfo) TableName() string {
	return "video_process_info"
}

func (vi *VideoProcessInfo) Mask(isAdmin bool) {
	vi.GenUID(common.DBTypeVideoProcessInfo)
}
