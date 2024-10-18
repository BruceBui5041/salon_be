package commentmodel

import (
	"video_server/common"
)

type CreateComment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32 `json:"user_id"`
	CourseID        string `json:"course_id"`
	Rate            uint8  `json:"rate"`
	Content         string `json:"content"`
}

func (cc *CreateComment) Mask(isAdmin bool) {
	cc.GenUID(common.DBTypeComment)
}
