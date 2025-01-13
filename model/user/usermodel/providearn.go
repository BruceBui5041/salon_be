package usermodel

import (
	"github.com/shopspring/decimal"
)

type EarningsSummary struct {
	TotalEarnings     decimal.Decimal  `json:"total_earnings"`
	CompletedBookings int              `json:"completed_bookings"`
	TotalHours        float64          `json:"total_hours"`
	Period            string           `json:"period"`
	MonthlyBreakdown  []MonthlyEarning `json:"monthly_breakdown,omitempty"`
}

type MonthlyEarning struct {
	Month             string          `json:"month"`
	Earnings          decimal.Decimal `json:"earnings"`
	CompletedBookings int             `json:"completed_bookings"`
	Hours             float64         `json:"hours"`
}

type GetEarningsRequest struct {
	Year  int `form:"year" binding:"required"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}
