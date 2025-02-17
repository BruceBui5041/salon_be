package groupprovidermodel

import "mime/multipart"

type GroupProviderCreateRequest struct {
	JSON   string                  `form:"json"`
	Images []*multipart.FileHeader `form:"images"`
}

type GroupProviderCreate struct {
	RequesterID uint32                  `json:"-"`
	OwnerStrID  string                  `json:"owner_id"`
	Name        string                  `json:"name" binding:"required"`
	Code        string                  `json:"code" binding:"required"`
	Description string                  `json:"description"`
	Images      []*multipart.FileHeader `json:"-"`
}
