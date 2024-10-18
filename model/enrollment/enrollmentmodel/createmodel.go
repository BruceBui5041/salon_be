package enrollmentmodel

import (
	"video_server/common"
)

type CreateEnrollment struct {
	common.SQLModel `json:",inline"`
	UserID          string `json:"user_id" form:"user_id"`
	CourseID        string `json:"course_id" form:"course_id"`
	PaymentID       string `json:"payment_id" form:"payment_id"`
}

func (ce *CreateEnrollment) Mask(isAdmin bool) {
	ce.GenUID(common.DBTypeEnrollment)
}
