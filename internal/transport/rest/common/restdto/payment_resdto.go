package restdto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentCreateRestDTO struct {
	RequestedAmount float64 `json:"requested_amount" binding:"required"`
	Email           *string `json:"email"`
}

type PaymentCreateResponseRestDTO struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
