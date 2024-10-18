package commentmodel

import (
	"salon_be/common"
)

type UpdateComment struct {
	common.SQLModel `json:",inline"`
	Content         string `json:"content" form:"content"`
	Rate            uint8  `json:"rate" form:"rate"`
}

func (uc *UpdateComment) Mask(isAdmin bool) {
	uc.GenUID(common.DBTypeComment)
}
