package categorymodel

import (
	"salon_be/common"
)

type CreateCategory struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name" gorm:"not null;size:100"`
	Description     string `json:"description"`
}

func (cc *CreateCategory) Mask(isAdmin bool) {
	cc.GenUID(common.DBTypeCategory)
}
