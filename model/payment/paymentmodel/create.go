package paymentmodel

import (
	"video_server/common"
)

type CreatePayment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32   `json:"user_id" form:"-"`
	Amount          float64  `json:"amount" form:"amount"`
	Currency        string   `json:"currency" form:"currency"`
	PaymentMethod   string   `json:"payment_method" form:"payment_method"`
	TransactionID   string   `json:"transaction_id" form:"transaction_id"`
	CourseIDs       []string `json:"course_ids" form:"course_ids"`
}

func (cp *CreatePayment) Mask(isAdmin bool) {
	cp.GenUID(common.DBTypePayment)
}
