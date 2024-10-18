package commentmodel

import (
	"salon_be/common"
	models "salon_be/model"
)

type CommentRes struct {
	common.SQLModel `json:",inline"`
	Content         string         `json:"content"`
	Rate            uint8          `json:"rate"`
	User            *models.User   `json:"user,omitempty"`
	Course          *models.Course `json:"course,omitempty"`
}