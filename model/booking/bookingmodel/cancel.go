package bookingmodel

type CancelBooking struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
	UserID             uint32 `json:"-"`
	IsUserRole         bool   `json:"-"`
}
