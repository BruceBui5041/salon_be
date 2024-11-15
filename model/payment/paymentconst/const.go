package paymentconst

const (
	PaymentMethodCash = "cash"
	PaymentMethodCard = "card"
)

const (
	TransactionStatusPending   = "pending"
	TransactionStatusCancelled = "cancelled"
	TransactionStatusCompleted = "completed"
)

var (
	PaymentMethods = []string{PaymentMethodCash, PaymentMethodCash}
)
