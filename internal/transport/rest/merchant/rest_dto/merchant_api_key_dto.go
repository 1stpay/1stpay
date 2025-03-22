package restdto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAPIKeyRequestDTO struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateAPIKeyResponseDTO struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	APIKey    string     `json:"api_key"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  bool       `json:"is_active"`
}
