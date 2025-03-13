package model

import (
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Blockchain struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"not null"`
	Logo      *string
	IsActive  bool             `gorm:"not null;default:false"`
	ChainType enum.NetworkType `gorm:"type:network_type;not null"`
	Config    datatypes.JSON
}
