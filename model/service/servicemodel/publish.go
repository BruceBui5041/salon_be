package servicemodel

import "salon_be/common"

type PublishServiceRequest struct {
	ServiceID        string `json:"id" form:"id"`
	ServiceVersionID string `json:"service_version_id" form:"service_version_id"`
}

func (ps *PublishServiceRequest) GetServiceLocalId() (uint32, error) {
	serviceUID, err := common.FromBase58(ps.ServiceID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}
	return serviceUID.GetLocalID(), nil
}

func (ps *PublishServiceRequest) GetServiceVersionLocalId() (uint32, error) {
	versionUID, err := common.FromBase58(ps.ServiceVersionID)
	if err != nil {
		return 0, common.ErrInvalidRequest(err)
	}
	return versionUID.GetLocalID(), nil
}
