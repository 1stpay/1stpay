package enum

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusNotFilled PaymentStatus = "not_filled"
)

type PaymentAMLStatus string

const (
	PaymentAMLStatusPassed  PaymentAMLStatus = "passed"
	PaymentAMLStatusFailed  PaymentAMLStatus = "failed"
	PaymentAMLStatusPending PaymentAMLStatus = "pending"
)
