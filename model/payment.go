package models

import (
	"salon_be/common"
	"salon_be/model/payment/paymentconst"

	"gorm.io/gorm"
)

const PaymentEntityName = "Payment"

type Payment struct {
	common.SQLModel   `json:",inline"`
	UserID            uint32  `json:"user_id" gorm:"column:user_id;not null"`
	Amount            float64 `json:"amount" gorm:"column:amount;type:decimal(10,2);not null"`
	Currency          string  `json:"currency" gorm:"column:currency;type:varchar(3);not null"`
	PaymentMethod     string  `json:"payment_method" gorm:"column:payment_method;type:varchar(50);not null"`
	TransactionID     string  `json:"transaction_id" gorm:"column:transaction_id;type:varchar(100);uniqueIndex"`
	TransactionStatus string  `json:"transaction_status" gorm:"column:transaction_status;type:varchar(50)"`
	User              User    `json:"user" gorm:"foreignKey:UserID"`
	// Enrollments       []Enrollment `json:"enrollments" gorm:"foreignKey:PaymentID"`
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:PaymentID"`
}

func (Payment) TableName() string {
	return "payment"
}

func (p *Payment) Mask(isAdmin bool) {
	p.GenUID(common.DBTypePayment)
}

func (p *Payment) AfterFind(tx *gorm.DB) (err error) {
	p.Mask(false)
	return
}

// BeforeCreate hook to calculate discounted price
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.TransactionStatus == "" {
		p.TransactionStatus = paymentconst.TransactionStatusPending
	}

	return nil
}
