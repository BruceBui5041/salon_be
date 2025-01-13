package userrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/user/usermodel"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BookingStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) ([]models.Booking, error)
}

type PaymentStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.Payment, error)
}

type providerEarningsRepo struct {
	db           *gorm.DB
	bookingStore BookingStore
	paymentStore PaymentStore
}

func NewProviderEarningsRepo(
	db *gorm.DB,
	bookingStore BookingStore,
	paymentStore PaymentStore,
) *providerEarningsRepo {
	return &providerEarningsRepo{
		db:           db,
		bookingStore: bookingStore,
		paymentStore: paymentStore,
	}
}

func (r *providerEarningsRepo) CalculateEarnings(
	ctx context.Context,
	providerID uint32,
	fromDate, toDate time.Time,
) (*usermodel.EarningsSummary, error) {
	var summary usermodel.EarningsSummary

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// First get the provider's completed bookings
		providerBookings := tx.Model(&models.Booking{}).
			Select("payment_id").
			Where("service_man_id = ?", providerID).
			Where("booking_status = ?", models.BookingStatusCompleted).
			Where("completed_at BETWEEN ? AND ?", fromDate, toDate)

		// Base query using the provider's bookings to check payments
		baseQuery := tx.Model(&models.Booking{}).
			Where("service_man_id = ?", providerID).
			Where("booking_status = ?", models.BookingStatusCompleted).
			Where("completed_at BETWEEN ? AND ?", fromDate, toDate).
			Where("payment_id IN (?)",
				tx.Model(&models.Payment{}).
					Select("id").
					Where("id IN (?)", providerBookings).
					Where("transaction_status = ?", "completed"),
			)

		// Calculate totals
		var totalResult struct {
			TotalEarnings     decimal.Decimal
			CompletedBookings int
			TotalHours        float64
		}

		err := baseQuery.Select(`
			COALESCE(SUM(CASE 
				WHEN discounted_price IS NOT NULL THEN discounted_price 
				ELSE price 
			END), 0) as total_earnings,
			COUNT(*) as completed_bookings,
			COALESCE(SUM(duration), 0) / 60.0 as total_hours
		`).Scan(&totalResult).Error

		if err != nil {
			return common.ErrDB(err)
		}

		summary.TotalEarnings = totalResult.TotalEarnings
		summary.CompletedBookings = totalResult.CompletedBookings
		summary.TotalHours = totalResult.TotalHours

		// Get monthly breakdown
		var monthlyResults []struct {
			Month             string
			Earnings          decimal.Decimal
			CompletedBookings int
			Hours             float64
		}

		err = baseQuery.Select(`
			DATE_FORMAT(completed_at, '%Y-%m') as month,
			COALESCE(SUM(CASE 
				WHEN discounted_price IS NOT NULL THEN discounted_price 
				ELSE price 
			END), 0) as earnings,
			COUNT(*) as completed_bookings,
			COALESCE(SUM(duration), 0) / 60.0 as hours
		`).
			Group("DATE_FORMAT(completed_at, '%Y-%m')").
			Order("month ASC").
			Scan(&monthlyResults).Error

		if err != nil {
			return common.ErrDB(err)
		}

		summary.MonthlyBreakdown = make([]usermodel.MonthlyEarning, len(monthlyResults))
		for i, result := range monthlyResults {
			summary.MonthlyBreakdown[i] = usermodel.MonthlyEarning{
				Month:             result.Month,
				Earnings:          result.Earnings,
				CompletedBookings: result.CompletedBookings,
				Hours:             result.Hours,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	summary.Period = fromDate.Format("2006-01")
	return &summary, nil
}
