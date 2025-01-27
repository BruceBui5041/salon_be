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
	summary := usermodel.EarningsSummary{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get bookings count by status
		var statusCounts struct {
			PendingBookings   int
			CancelledBookings int
			ConfirmedBookings int
		}

		err := tx.Model(&models.Booking{}).
			Select(`
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as pending_bookings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as cancelled_bookings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as confirmed_bookings
			`, models.BookingStatusPending, models.BookingStatusCancelled, models.BookingStatusConfirmed).
			Where("service_man_id = ?", providerID).
			Where("booking_date BETWEEN ? AND ?", fromDate, toDate).
			Scan(&statusCounts).Error

		if err != nil {
			return common.ErrDB(err)
		}

		summary.PendingBookings = statusCounts.PendingBookings
		summary.CancelledBookings = statusCounts.CancelledBookings
		summary.ConfirmedBookings = statusCounts.ConfirmedBookings

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
			TotalCommission   decimal.Decimal
		}

		err = baseQuery.Select(`
			COALESCE(SUM(CASE 
				WHEN discounted_price IS NOT NULL THEN discounted_price 
				ELSE price 
			END), 0) as total_earnings,
			COUNT(*) as completed_bookings,
			COALESCE(SUM(duration), 0) / 60.0 as total_hours,
			COALESCE(SUM(CASE 
				WHEN discounted_price IS NOT NULL THEN discounted_price * COALESCE(commission.percentage, 0) / 100 
				ELSE price * COALESCE(commission.percentage, 0) / 100 
			END), 0) as total_commission
		`).
			Joins("LEFT JOIN commission ON booking.commission_id = commission.id").
			Scan(&totalResult).Error

		if err != nil {
			return common.ErrDB(err)
		}

		summary.TotalEarnings = totalResult.TotalEarnings.Sub(totalResult.TotalCommission)
		summary.CompletedBookings = totalResult.CompletedBookings
		summary.TotalHours = totalResult.TotalHours
		summary.TotalCommission = totalResult.TotalCommission

		// Get monthly breakdown with status counts
		var monthlyResults []struct {
			Month             string
			Earnings          decimal.Decimal
			CompletedBookings int
			PendingBookings   int
			CancelledBookings int
			ConfirmedBookings int
			Hours             float64
			Commission        decimal.Decimal
		}

		err = tx.Model(&models.Booking{}).
			Select(`
				DATE_FORMAT(booking_date, '%Y-%m') as month,
				COALESCE(SUM(CASE 
					WHEN booking_status = ? AND discounted_price IS NOT NULL THEN discounted_price 
					WHEN booking_status = ? THEN price 
					ELSE 0 
				END), 0) as earnings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as completed_bookings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as pending_bookings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as cancelled_bookings,
				COUNT(CASE WHEN booking_status = ? THEN 1 END) as confirmed_bookings,
				COALESCE(SUM(CASE WHEN booking_status = ? THEN duration ELSE 0 END), 0) / 60.0 as hours,
				COALESCE(SUM(CASE 
					WHEN booking_status = ? AND discounted_price IS NOT NULL THEN discounted_price * COALESCE(commission.percentage, 0) / 100 
					WHEN booking_status = ? THEN price * COALESCE(commission.percentage, 0) / 100 
					ELSE 0 
				END), 0) as commission
			`,
				models.BookingStatusCompleted, models.BookingStatusCompleted,
				models.BookingStatusCompleted,
				models.BookingStatusPending,
				models.BookingStatusCancelled,
				models.BookingStatusConfirmed,
				models.BookingStatusCompleted,
				models.BookingStatusCompleted, models.BookingStatusCompleted).
			Joins("LEFT JOIN commission ON booking.commission_id = commission.id").
			Where("service_man_id = ?", providerID).
			Where("booking_date BETWEEN ? AND ?", fromDate, toDate).
			Group("DATE_FORMAT(booking_date, '%Y-%m')").
			Order("month ASC").
			Scan(&monthlyResults).Error

		if err != nil {
			return common.ErrDB(err)
		}

		summary.MonthlyBreakdown = make([]usermodel.MonthlyEarning, len(monthlyResults))
		for i, result := range monthlyResults {
			summary.MonthlyBreakdown[i] = usermodel.MonthlyEarning{
				Month:             result.Month,
				Earnings:          result.Earnings.Sub(result.Commission),
				CompletedBookings: result.CompletedBookings,
				PendingBookings:   result.PendingBookings,
				CancelledBookings: result.CancelledBookings,
				ConfirmedBookings: result.ConfirmedBookings,
				Hours:             result.Hours,
				Commission:        result.Commission,
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
