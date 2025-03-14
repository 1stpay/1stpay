package restdto

type InvoiceCreateRestDTO struct {
	RequestedAmount float64 `json:"requested_amount" binding:"required"`
	Email           *string `json:"email"`
}
