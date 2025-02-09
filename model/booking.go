package models

import (
	"errors"
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	modelhelper.RegisterModel(Booking{})
}

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
	UserID             uint32            `json:"-" gorm:"column:user_id;not null;index"`
	User               *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ServiceVersions    []*ServiceVersion `json:"service_versions,omitempty" gorm:"many2many:m2mbooking_service_version;foreignKey:Id;joinForeignKey:BookingID;References:Id;joinReferences:ServiceVersionID"`
	ServiceManID       uint32            `json:"-" gorm:"column:service_man_id;not null;index"`
	ServiceMan         *User             `json:"service_man,omitempty" gorm:"foreignKey:ServiceManID"`
	PaymentID          *uint32           `json:"-" gorm:"column:payment_id;index"`
	Payment            *Payment          `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	CouponID           *uint32           `json:"-" gorm:"column:coupon_id;index"`
	Coupon             *Coupon           `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	BookingStatus      BookingStatus     `json:"booking_status" gorm:"column:booking_status;type:ENUM('pending','confirmed','in_progress','completed','cancelled','no_show');default:pending"`
	ConfirmedDate      *time.Time        `json:"confirmed_date,omitempty" gorm:"column:confirmed_date;type:datetime"`
	BookingDate        time.Time         `json:"booking_date" gorm:"column:booking_date;type:datetime;not null"`
	Duration           uint32            `json:"duration" gorm:"column:duration;not null"`
	Price              decimal.Decimal   `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	DiscountedPrice    *decimal.Decimal  `json:"discounted_price" gorm:"column:discounted_price;type:decimal(10,2)"`
	DiscountAmount     decimal.Decimal   `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2)"`
	Notes              string            `json:"notes" gorm:"column:notes;type:text"`
	CancelledByID      *uint32           `json:"-" gorm:"column:cancelled_by_id;index"`
	CancelledBy        *User             `json:"cancelled_by,omitempty" gorm:"foreignKey:CancelledByID"`
	CancellationReason string            `json:"cancellation_reason,omitempty" gorm:"column:cancellation_reason;type:text"`
	CancelledAt        *time.Time        `json:"cancelled_at,omitempty" gorm:"column:cancelled_at;type:datetime"`
	CompletedAt        *time.Time        `json:"completed_at,omitempty" gorm:"column:completed_at;type:datetime"`
	Notifications      []*Notification   `json:"notifications,omitempty" gorm:"foreignKey:BookingID"`
	CommissionID       *uint32           `json:"-" gorm:"column:commission_id;index"`
	Commission         *Commission       `json:"commission,omitempty" gorm:"foreignKey:CommissionID"`

	// BookingLocation   *BookingLocation    `json:"booking_location,omitempty" gorm:"foreignKey:BookingID"`
	// Locations         []*Location         `json:"locations,omitempty" gorm:"foreignKey:BookingID"`
	// DistanceTrackings []*DistanceTracking `json:"distance_trackings,omitempty" gorm:"foreignKey:BookingID"`
}

func (Booking) TableName() string {
	return "booking"
}

func (b *Booking) Mask(isAdmin bool) {
	b.GenUID(common.DBTypeBooking)
}

func (b *Booking) AfterFind(tx *gorm.DB) (err error) {
	b.Mask(false)
	return
}

func (b *Booking) CalculateDiscountedPrice() error {
	if len(b.ServiceVersions) == 0 {
		return errors.New("at least one service version is required")
	}

	// Calculate total price from all service versions
	totalPrice := decimal.Zero
	for _, sv := range b.ServiceVersions {
		totalPrice = totalPrice.Add(sv.Price)
	}
	b.Price = totalPrice

	// Calculate total discounted price if any service versions have discounts
	// one booking can have multiple service versions with different discounts
	totalDiscountedPrice := decimal.Zero
	hasDiscounts := false
	for _, sv := range b.ServiceVersions {
		if sv.DiscountedPrice != nil {
			hasDiscounts = true
			totalDiscountedPrice = totalDiscountedPrice.Add(sv.DiscountedPrice.Decimal)
		} else {
			totalDiscountedPrice = totalDiscountedPrice.Add(sv.Price)
		}
	}
	if hasDiscounts {
		b.DiscountedPrice = &totalDiscountedPrice
	}

	if b.Coupon == nil {
		return nil
	}

	if err := b.Coupon.IsValid(totalDiscountedPrice); err != nil {
		return err
	}

	discountAmount := b.Coupon.CalculateDiscount(totalDiscountedPrice)
	b.DiscountAmount = discountAmount

	finalPrice := totalDiscountedPrice.Sub(discountAmount)
	b.DiscountedPrice = &finalPrice

	return nil
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.BookingStatus == "" {
		b.BookingStatus = BookingStatusPending
	}

	// Calculate discounted price if coupon is provided
	if err := b.CalculateDiscountedPrice(); err != nil {
		return err
	}

	// Fetch commission based on service man's role
	if b.ServiceManID != 0 {
		var serviceMan User
		if err := tx.Preload("Roles").First(&serviceMan, b.ServiceManID).Error; err != nil {
			return err
		}

		// Filter roles that contain "PROVIDER"
		providerRoles := lo.Filter(serviceMan.Roles, func(role *Role, _ int) bool {
			return strings.Contains(role.Code, "PROVIDER")
		})

		if len(providerRoles) > 0 {
			// Get all commissions for provider roles
			var commissions []Commission
			if err := tx.Where(
				"role_id IN (?) AND status = ?",
				lo.Map(providerRoles, func(r *Role, _ int) uint32 { return r.Id }),
				common.StatusActive,
			).Find(&commissions).Error; err != nil {
				return err
			}

			// Find commission with lowest percentage
			if len(commissions) > 0 {
				lowestCommission := lo.MinBy(commissions, func(a, b Commission) bool {
					return a.Percentage < b.Percentage
				})
				b.CommissionID = &lowestCommission.Id
			}
		}
	}

	return nil
}
