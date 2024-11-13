package bookingmodel

type CompleteBooking struct {
	UserID     uint32 `json:"-"`
	IsUserRole bool   `json:"-"`
}
