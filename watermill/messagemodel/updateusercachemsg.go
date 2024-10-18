package messagemodel

type EnrollmentChangeInfo struct {
	UserId            string `json:"user_id"`
	ServiceId         string `json:"service_id"`
	ServiceSlug       string `json:"service_slug"`
	EnrollmentId      string `json:"enrollment_id"`
	PaymentId         string `json:"payment_id"`
	TransactionStatus string `json:"transaction_status"`
}
