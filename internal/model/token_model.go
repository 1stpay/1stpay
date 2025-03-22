package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Token struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string    `gorm:"not null"`
	Symbol          string    `gorm:"not null"`
	ContractAddress string    `gorm:"not null"`
	Decimals        int       `gorm:"not null"`
	Logo            *string
	BlockchainID    uuid.UUID  `gorm:"type:uuid;not null"`
	Blockchain      Blockchain `gorm:"foreignKey:BlockchainID"`
	IsNative        bool       `gorm:"not null;default:false"`
	IsActive        bool       `gorm:"not null;default:false"`
	Config          datatypes.JSON
}
