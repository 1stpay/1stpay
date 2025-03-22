package restdto

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type TokenListResponseDTO struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	Symbol       string      `json:"symbol"`
	BlockchainID uuid.UUID   `json:"blockchain_id"`
	Logo         null.String `json:"logo"`
	IsNative     bool        `json:"is_native"`
	IsActive     bool        `json:"is_active"`
}
