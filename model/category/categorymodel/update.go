package categorymodel

import "mime/multipart"

type UpdateCategory struct {
	Name        *string               `json:"name" form:"name"`
	Description *string               `json:"description" form:"description"`
	Code        *string               `json:"code" form:"code"`
	Image       *multipart.FileHeader `json:"image" form:"image"`
}
