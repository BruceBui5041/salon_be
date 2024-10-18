package progressmodel

import (
	"time"
	"video_server/common"
)

type CreateProgress struct {
	common.SQLModel `json:",inline"`
	UserID          uint32    `json:"user_id" gorm:"column:user_id;index"`
	VideoID         string    `json:"video_id" gorm:"-"`
	WatchedSeconds  uint32    `json:"watched_seconds" gorm:"column:watched_seconds;default:0"`
	Completed       bool      `json:"completed" gorm:"column:completed;default:false"`
	LastWatched     time.Time `json:"last_watched" gorm:"column:last_watched;autoUpdateTime"`
}

func (cp *CreateProgress) Mask(isAdmin bool) {
	cp.GenUID(common.DBTypeProgress)
}
