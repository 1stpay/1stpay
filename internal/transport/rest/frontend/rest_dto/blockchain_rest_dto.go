package restdto

import (
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type BlockchainListResponseDTO struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Logo      null.String      `json:"logo"`
	IsActive  bool             `json:"is_active"`
	ChainType enum.NetworkType `json:"chain_type"`
}
