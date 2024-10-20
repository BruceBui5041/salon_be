package categorymodel

import (
	"mime/multipart"
	"salon_be/common"
)

type CreateCategory struct {
	common.SQLModel `json:",inline"`
	Name            string                `json:"name" form:"name"`
	Description     string                `json:"description" form:"description"`
	Code            string                `json:"code" form:"code"`
	Image           *multipart.FileHeader `json:"image" form:"image"`
}

func (cc *CreateCategory) Mask(isAdmin bool) {
	cc.GenUID(common.DBTypeCategory)
}
