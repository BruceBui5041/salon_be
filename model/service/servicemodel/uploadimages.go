package servicemodel

import (
	"mime/multipart"
	"salon_be/common"
)

type UploadImages struct {
	UploadedBy       uint32                  `json:"uploaded_by" form:"uploaded_by"`
	ServiceID        string                  `json:"service_id" form:"service_id"`
	ServiceVersionID *string                 `json:"service_version_id" form:"service_version_id"`
	Images           []*multipart.FileHeader `json:"images" form:"images"`
}

func (ui *UploadImages) GetServiceIDLocalId() (uint32, error) {
	serviceUID, err := common.FromBase58(ui.ServiceID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return serviceUID.GetLocalID(), nil
}

func (ui *UploadImages) GetServiceVersionIDLocalId() (uint32, error) {
	serviceVersionUID, err := common.FromBase58(*ui.ServiceVersionID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}

	return serviceVersionUID.GetLocalID(), nil
}
