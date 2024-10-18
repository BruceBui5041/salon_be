package commentmodel

import (
	"video_server/common"
	models "video_server/model"
)

type CommentRes struct {
	common.SQLModel `json:",inline"`
	Content         string         `json:"content"`
	Rate            uint8          `json:"rate"`
	User            *models.User   `json:"user,omitempty"`
	Course          *models.Course `json:"course,omitempty"`
}
