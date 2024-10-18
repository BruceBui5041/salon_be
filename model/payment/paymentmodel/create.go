package paymentmodel

import (
	"salon_be/common"
)

type CreatePayment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32   `json:"user_id" form:"-"`
	Amount          float64  `json:"amount" form:"amount"`
	Currency        string   `json:"currency" form:"currency"`
	PaymentMethod   string   `json:"payment_method" form:"payment_method"`
	TransactionID   string   `json:"transaction_id" form:"transaction_id"`
	ServiceIDs      []string `json:"service_ids" form:"service_ids"`
}

func (cp *CreatePayment) Mask(isAdmin bool) {
	cp.GenUID(common.DBTypePayment)
}
