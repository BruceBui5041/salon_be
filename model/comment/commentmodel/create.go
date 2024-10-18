package commentmodel

import (
	"salon_be/common"
)

type CreateComment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32 `json:"user_id"`
	ServiceID       string `json:"service_id"`
	Rate            uint8  `json:"rate"`
	Content         string `json:"content"`
}

func (cc *CreateComment) Mask(isAdmin bool) {
	cc.GenUID(common.DBTypeComment)
}
