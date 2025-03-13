package restdto

import (
	"time"

	"github.com/google/uuid"
)

type MerchantCreateRequestDTO struct {
	Name string `json:"name" binding:"required"`
}

type MerchantCreateResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CommissionRate float64   `json:"commision_rate"`
}
type MerchantDetailResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CommissionRate float64   `json:"commision_rate"`
}

type MerchantTokenCreateRequestDTO struct {
	TokenID uuid.UUID `json:"token_id" binding:"required"`
	Active  bool      `json:"active" binding:"required"`
}

type MerchantTokenCreateResponseDTO struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	MerchantID uuid.UUID `json:"merchant_id" binding:"required"`
	TokenID    uuid.UUID `json:"token_id" binding:"required"`
	Active     bool      `json:"active" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"required"`
}
