package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt      time.Time `gorm:"not null;default:now()"`
	UpdatedAt      time.Time `gorm:"not null;default:now()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	User           User      `gorm:"foreignKey:UserID"`
	Name           string    `gorm:"not null"`
	CommissionRate float64   `gorm:"type:numeric(5,2);not null;default:0"`
}

type MerchantToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null"`
	Merchant   Merchant  `gorm:"foreignKey:MerchantID"`
	TokenID    uuid.UUID `gorm:"type:uuid;not null"`
	Token      Token     `gorm:"foreignKey:TokenID"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
}
