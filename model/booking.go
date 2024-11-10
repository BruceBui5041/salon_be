package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	BookingEntityName = "Booking"
)

type BookingStatus string

const (
	BookingStatusPending    BookingStatus = "pending"
	BookingStatusConfirmed  BookingStatus = "confirmed"
	BookingStatusInProgress BookingStatus = "in_progress"
	BookingStatusCompleted  BookingStatus = "completed"
	BookingStatusCancelled  BookingStatus = "cancelled"
	BookingStatusNoShow     BookingStatus = "no_show"
)

// Updated Booking model
type Booking struct {
	common.SQLModel    `json:",inline"`
	UserID             uint32           `json:"-" gorm:"column:user_id;not null;index"`
	User               *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ServiceVersionID   uint32           `json:"-" gorm:"column:service_version_id;not null;index"`
	ServiceVersion     *ServiceVersion  `json:"service_version,omitempty" gorm:"foreignKey:ServiceVersionID"`
	ServiceManID       uint32           `json:"-" gorm:"column:service_man_id;not null;index"`
	ServiceMan         *User            `json:"service_man,omitempty" gorm:"foreignKey:ServiceManID"`
	PaymentID          *uint32          `json:"-" gorm:"column:payment_id;index"`
	Payment            *Payment         `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	CouponID           *uint32          `json:"-" gorm:"column:coupon_id;index"`
	Coupon             *Coupon          `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	BookingStatus      BookingStatus    `json:"booking_status" gorm:"column:booking_status;type:ENUM('pending','confirmed','in_progress','completed','cancelled','no_show');default:pending"`
	BookingDate        time.Time        `json:"booking_date" gorm:"column:booking_date;type:datetime;not null"`
	Duration           uint32           `json:"duration" gorm:"column:duration;not null"`
	Price              decimal.Decimal  `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	DiscountedPrice    *decimal.Decimal `json:"discounted_price" gorm:"column:discounted_price;type:decimal(10,2)"`
	DiscountAmount     decimal.Decimal  `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2)"`
	Notes              string           `json:"notes" gorm:"column:notes;type:text"`
	CancellationReason string           `json:"cancellation_reason,omitempty" gorm:"column:cancellation_reason;type:text"`
	CancelledAt        *time.Time       `json:"cancelled_at,omitempty" gorm:"column:cancelled_at;type:datetime"`
	CompletedAt        *time.Time       `json:"completed_at,omitempty" gorm:"column:completed_at;type:datetime"`
}

func (Booking) TableName() string {
	return "booking"
}

func (b *Booking) CalculateDiscountedPrice() error {
	if b.Coupon == nil {
		return nil
	}

	// Validate coupon
	if err := b.Coupon.IsValid(b.Price); err != nil {
		return err
	}

	// Calculate discount amount
	discountAmount := b.Coupon.CalculateDiscount(b.Price)
	b.DiscountAmount = discountAmount

	// Calculate final price
	finalPrice := b.Price.Sub(discountAmount)
	b.DiscountedPrice = &finalPrice

	return nil
}

// BeforeCreate hook to calculate discounted price
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.BookingStatus == "" {
		b.BookingStatus = BookingStatusPending
	}

	// Calculate discounted price if coupon is provided
	if err := b.CalculateDiscountedPrice(); err != nil {
		return err
	}

	return nil
}

func init() {
	modelhelper.RegisterModel(Booking{})
}
