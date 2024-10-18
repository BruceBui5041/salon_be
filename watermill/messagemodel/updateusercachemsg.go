package messagemodel

type EnrollmentChangeInfo struct {
	UserId            string `json:"user_id"`
	CourseId          string `json:"course_id"`
	CourseSlug        string `json:"course_slug"`
	EnrollmentId      string `json:"enrollment_id"`
	PaymentId         string `json:"payment_id"`
	TransactionStatus string `json:"transaction_status"`
}
