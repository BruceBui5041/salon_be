package couponmodel

import (
	"mime/multipart"
	models "salon_be/model"
	"time"

	"github.com/shopspring/decimal"
)

type UpdateCouponRequest struct {
	JSON  string                `form:"json"`
	Image *multipart.FileHeader `form:"image"`
}

type UpdateCoupon struct {
	Status        *string               `json:"status"`
	Description   string                `json:"description"`
	DiscountType  models.DiscountType   `json:"discount_type"`
	DiscountValue decimal.Decimal       `json:"discount_value"`
	MinSpend      decimal.Decimal       `json:"min_spend"`
	MaxDiscount   decimal.Decimal       `json:"max_discount"`
	StartDate     time.Time             `json:"start_date"`
	EndDate       time.Time             `json:"end_date"`
	CreatorID     uint32                `json:"-"`
	UsageLimit    *int                  `json:"usage_limit"`
	Image         *multipart.FileHeader `json:"-"`
}
