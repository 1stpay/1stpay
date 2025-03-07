package enum

// PaymentStatus – тип для статуса платежа
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusNotFilled PaymentStatus = "not_filled"
)

// PaymentAMLStatus – тип для статуса AML проверки платежа
type PaymentAMLStatus string

const (
	PaymentAMLStatusPassed  PaymentAMLStatus = "passed"
	PaymentAMLStatusFailed  PaymentAMLStatus = "failed"
	PaymentAMLStatusPending PaymentAMLStatus = "pending"
)
