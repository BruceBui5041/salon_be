package models

import (
	"errors"
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/model/coupon/couponerror"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	modelhelper.RegisterModel(Coupon{})
}

const (
	CouponEntityName = "Coupon"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixedPrice DiscountType = "fixed_price"
)

// Coupon model
type Coupon struct {
	common.SQLModel `json:",inline"`
	Code            string          `json:"code" gorm:"column:code;uniqueIndex;not null;type:varchar(20)"`
	Description     string          `json:"description" gorm:"column:description;type:text"`
	DiscountType    DiscountType    `json:"discount_type" gorm:"column:discount_type;type:ENUM('percentage','fixed_price');not null"`
	DiscountValue   decimal.Decimal `json:"discount_value" gorm:"column:discount_value;type:decimal(10,2);not null"`
	MinSpend        decimal.Decimal `json:"min_spend" gorm:"column:min_spend;type:decimal(10,2)"`
	MaxDiscount     decimal.Decimal `json:"max_discount" gorm:"column:max_discount;type:decimal(10,2)"`
	StartDate       time.Time       `json:"start_date" gorm:"column:start_date;type:datetime;not null"`
	EndDate         time.Time       `json:"end_date" gorm:"column:end_date;type:datetime;not null"`
	UsageLimit      *int            `json:"usage_limit" gorm:"column:usage_limit"`
	UsageCount      int             `json:"usage_count" gorm:"column:usage_count;default:0"`
	Bookings        []*Booking      `json:"bookings,omitempty" gorm:"foreignKey:CouponID"`
	CreatorID       uint32          `json:"-" gorm:"column:creator_id;not null;index"`
	Creator         *User           `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	ImageID         *uint32         `json:"-" gorm:"column:image_id"`
	Image           *Image          `json:"image,omitempty" gorm:"foreignKey:ImageID"`
	CategoryID      *uint32         `json:"-" gorm:"column:category_id;foreignKey:Id;constraint:OnDelete:SET NULL;"`
	Category        *Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Coupon) TableName() string {
	return "coupon"
}

func (c *Coupon) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCoupon)
}

func (c *Coupon) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	return
}

func (c *Coupon) BeforeCreate(tx *gorm.DB) error {
	c.Status = common.StatusInactive
	return nil
}

func (c *Coupon) BeforeUpdate(tx *gorm.DB) error {
	if c.Status == common.StatusSuspended {
		return nil
	}

	var existingRecord Coupon
	if err := tx.Model(&Coupon{}).
		Where("id = ?", c.Id).
		First(&existingRecord).
		Error; err != nil {
		return err
	}

	if existingRecord.Status == common.StatusActive || existingRecord.Status == common.StatusSuspended {
		return couponerror.ErrCouponHasBeenActivated(errors.New("coupon has been activated"))
	}
	return nil
}

func (c *Coupon) IsValid(originalPrice decimal.Decimal) error {
	now := time.Now().UTC()

	if c.Status == common.StatusInactive {
		return errors.New("coupon is not active")
	}

	if c.Status == common.StatusSuspended {
		return errors.New("coupon is suspended")
	}

	if now.Before(c.StartDate) {
		return errors.New("coupon is not yet active")
	}

	if now.After(c.EndDate) {
		return errors.New("coupon has expired")
	}

	if c.UsageLimit != nil && *c.UsageLimit != 0 && c.UsageCount >= *c.UsageLimit {
		return errors.New("coupon usage limit exceeded")
	}

	if originalPrice.LessThan(c.MinSpend) {
		return errors.New("booking amount does not meet minimum spend requirement")
	}

	return nil
}

func (c *Coupon) CalculateDiscount(originalPrice decimal.Decimal) decimal.Decimal {
	if c.DiscountType == DiscountTypeFixedPrice {
		if originalPrice.LessThan(c.DiscountValue) {
			return originalPrice
		}
		return c.DiscountValue
	}

	// Percentage discount calculation
	discount := originalPrice.Mul(c.DiscountValue.Div(decimal.NewFromInt(100)))

	// Check if discount exceeds maximum discount allowed (for percentage type)
	if c.MaxDiscount.IsPositive() && discount.GreaterThan(c.MaxDiscount) {
		return c.MaxDiscount
	}

	return discount
}
