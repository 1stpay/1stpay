package restdto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentCreateResponseRestDTO struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
