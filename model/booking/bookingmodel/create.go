package bookingmodel

import (
	"time"
)

type CreateBooking struct {
	ServiceVersionID uint32    `json:"service_version_id" binding:"required"`
	ServiceManID     uint32    `json:"service_man_id" binding:"required"`
	CouponID         *uint32   `json:"coupon_id"`
	BookingDate      time.Time `json:"booking_date" binding:"required"`
	Notes            string    `json:"notes"`
	UserID           uint32    `json:"-"`
}
