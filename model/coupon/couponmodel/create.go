package couponmodel

import (
	"mime/multipart"
	models "salon_be/model"
	"time"

	"github.com/shopspring/decimal"
)

type CreateCouponRequest struct {
	JSON  string                `form:"json"`
	Image *multipart.FileHeader `form:"image"`
}

type CreateCoupon struct {
	Code          string                `json:"code" binding:"required"`
	Description   string                `json:"description"`
	DiscountType  models.DiscountType   `json:"discount_type" binding:"required"`
	DiscountValue decimal.Decimal       `json:"discount_value" binding:"required"`
	MinSpend      decimal.Decimal       `json:"min_spend"`
	MaxDiscount   decimal.Decimal       `json:"max_discount"`
	StartDate     time.Time             `json:"start_date" binding:"required"`
	EndDate       time.Time             `json:"end_date" binding:"required"`
	UsageLimit    *int                  `json:"usage_limit"`
	CreatorID     uint32                `json:"-"`
	Image         *multipart.FileHeader `json:"-"`
}
