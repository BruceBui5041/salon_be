package enrollmentmodel

import (
	"salon_be/common"
)

type CreateEnrollment struct {
	common.SQLModel `json:",inline"`
	UserID          string `json:"user_id" form:"user_id"`
	ServiceID       string `json:"service_id" form:"service_id"`
	PaymentID       string `json:"payment_id" form:"payment_id"`
}

func (ce *CreateEnrollment) Mask(isAdmin bool) {
	ce.GenUID(common.DBTypeEnrollment)
}
