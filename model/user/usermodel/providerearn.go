package usermodel

import (
	"github.com/shopspring/decimal"
)

type EarningsSummary struct {
	TotalEarnings     decimal.Decimal  `json:"total_earnings"`
	CompletedBookings int              `json:"completed_bookings"`
	PendingBookings   int              `json:"pending_bookings"`
	CancelledBookings int              `json:"cancelled_bookings"`
	ConfirmedBookings int              `json:"confirmed_bookings"`
	TotalHours        float64          `json:"total_hours"`
	TotalCommission   decimal.Decimal  `json:"total_commission"`
	Period            string           `json:"period"`
	MonthlyBreakdown  []MonthlyEarning `json:"monthly_breakdown,omitempty"`
}

type MonthlyEarning struct {
	Month             string          `json:"month"`
	Earnings          decimal.Decimal `json:"earnings"`
	CompletedBookings int             `json:"completed_bookings"`
	PendingBookings   int             `json:"pending_bookings"`
	CancelledBookings int             `json:"cancelled_bookings"`
	ConfirmedBookings int             `json:"confirmed_bookings"`
	Hours             float64         `json:"hours"`
	Commission        decimal.Decimal `json:"commission"`
}

type GetEarningsRequest struct {
	Year  int `form:"year" binding:"required"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}
