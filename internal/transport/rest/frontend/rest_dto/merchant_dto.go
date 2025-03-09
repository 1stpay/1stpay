package restdto

import (
	"time"

	"github.com/google/uuid"
)

type MerchantCreateRequestDTO struct {
	Name string `json:"name"`
}

type MerchantCreateResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CommissionRate float64   `json:"commision_rate"`
}
