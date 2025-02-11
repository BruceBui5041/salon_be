package certificatemodel

import "mime/multipart"

type CreateCertificateInput struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Type      string                `form:"type" binding:"required"`
	OwnerID   uint32                `json:"owner_id"`
	CreatorID uint32                `json:"creator_id"`
}
