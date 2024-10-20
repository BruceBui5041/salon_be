package categorymodel

import (
	"salon_be/common"
)

type CreateCategory struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Code            string `json:"code"`
}

func (cc *CreateCategory) Mask(isAdmin bool) {
	cc.GenUID(common.DBTypeCategory)
}
