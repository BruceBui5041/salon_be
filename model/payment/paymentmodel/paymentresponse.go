package paymentmodel

import (
	"salon_be/common"
	models "salon_be/model"
)

type PaymentResponse struct {
	common.SQLModel   `json:",inline"`
	Amount            float64              `json:"amount"`
	Currency          string               `json:"currency"`
	PaymentMethod     string               `json:"payment_method"`
	TransactionID     string               `json:"transaction_id"`
	TransactionStatus string               `json:"transaction_status"`
	Enrollments       []EnrollmentResponse `json:"enrollments"`
}

func (p *PaymentResponse) Mask(isAdmin bool) {
	p.GenUID(common.DBTypePayment)
}

type ServiceResponse struct {
	common.SQLModel `json:",inline"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Creator         models.User     `json:"creator"`
	Category        models.Category `json:"category"`
	Slug            string          `json:"slug"`
	Thumbnail       string          `json:"thumbnail"`
	Price           uint64          `json:"price"`
	DiscountedPrice *uint64         `json:"discounted_price" `
	DifficultyLevel string          `json:"difficulty_level"`
}

func (c *ServiceResponse) Mask(isAdmin bool) {
	c.GenUID(common.DbTypeServiceVersion)
}

type EnrollmentResponse struct {
	common.SQLModel `json:",inline"`
	Service         ServiceResponse `json:"service"`
}

func (EnrollmentResponse) TableName() string {
	return "enrollment"
}

func (enr *EnrollmentResponse) Mask(isAdmin bool) {
	enr.GenUID(common.DBTypeEnrollment)
}
