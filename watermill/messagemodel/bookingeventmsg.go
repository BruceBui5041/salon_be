package messagemodel

const (
	BookingCreatedEvent   = "booking_created"
	BookingAcceptedEvent  = "booking_accepted"
	BookingCompletedEvent = "booking_completed"
	BookingCancelledEvent = "booking_cancelled"
)

type BookingEventMsg struct {
	BookingID uint32 `json:"booking_id"`
	Event     string `json:"event"`
}
