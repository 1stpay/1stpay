package model

import (
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/google/uuid"
)

type Payment struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt          time.Time  `gorm:"not null;default:now()"`
	UpdatedAt          time.Time  `gorm:"not null;default:now()"`
	TokenID            uuid.UUID  `gorm:"type:uuid;not null"`
	Token              Token      `gorm:"foreignKey:TokenID"`
	BlockchainID       uuid.UUID  `gorm:"type:uuid;not null"`
	Blockchain         Blockchain `gorm:"foreignKey:BlockchainID"`
	MerchantID         uuid.UUID  `gorm:"type:uuid;not null"`
	Merchant           Merchant   `gorm:"foreignKey:MerchantID"`
	RequestedAmount    float64    `gorm:"type:numeric(20,8);not null;default:0"`
	PaidAmount         float64    `gorm:"type:numeric(20,8);not null;default:0"`
	RequestedAmountWei int64      `gorm:"not null;default:0"`
	PaidAmountWei      int64      `gorm:"not null;default:0"`
	StableAmount       float64    `gorm:"type:numeric(20,8);not null"`
	CommissionAmount   float64    `gorm:"type:numeric(20,8);not null;default:0"`
	ExpireDate         *time.Time
	AMLStatus          enum.PaymentAMLStatus `gorm:"type:payment_aml_status"`
	Status             enum.PaymentStatus    `gorm:"type:payment_status;not null;default:'pending'"`
	InvoiceEmail       string
}

type PaymentAddress struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt    time.Time  `gorm:"not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()"`
	PaymentID    uuid.UUID  `gorm:"type:uuid;not null"`
	Payment      Payment    `gorm:"foreignKey:PaymentID"`
	BlockchainID uuid.UUID  `gorm:"type:uuid;not null"`
	Blockchain   Blockchain `gorm:"foreignKey:BlockchainID"`
	TokenID      uuid.UUID  `gorm:"type:uuid;not null"`
	Token        Token      `gorm:"foreignKey:TokenID"`
	PublicKey    string     `gorm:"not null"`
	PrivateKey   string     `gorm:"not null"`
}
