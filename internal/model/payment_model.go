package model

import (
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/google/uuid"
)

type Payment struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt        time.Time `gorm:"not null;default:now()"`
	UpdatedAt        time.Time `gorm:"not null;default:now()"`
	MerchantID       uuid.UUID `gorm:"type:uuid;not null"`
	Merchant         Merchant  `gorm:"foreignKey:MerchantID"`
	RequestedAmount  float64   `gorm:"type:numeric(20,8);not null;default:0"`
	PaidAmount       float64   `gorm:"type:numeric(20,8);not null;default:0"`
	CommissionAmount float64   `gorm:"type:numeric(20,8);not null;default:0"`
	ExpiresAt        *time.Time
	AMLStatus        *enum.PaymentAMLStatus `gorm:"type:payment_aml_status"`
	Status           enum.PaymentStatus     `gorm:"type:payment_status;not null;default:'pending'"`
	InvoiceEmail     *string
	UsedTokenID      *uuid.UUID `gorm:"type:uuid;"`
	UsedToken        *Merchant  `gorm:"foreignKey:UsedTokenID"`
}

type PaymentAddress struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt          time.Time `gorm:"not null;default:now()"`
	UpdatedAt          time.Time `gorm:"not null;default:now()"`
	PaymentID          uuid.UUID `gorm:"type:uuid;not null"`
	Payment            Payment   `gorm:"foreignKey:PaymentID"`
	TokenID            uuid.UUID `gorm:"type:uuid;not null"`
	Token              Token     `gorm:"foreignKey:TokenID"`
	PublicKey          string    `gorm:"not null"`
	PrivateKey         string    `gorm:"not null"`
	RequestedAmount    float64   `gorm:"type:numeric(20,8);not null;default:0"`
	PaidAmount         float64   `gorm:"type:numeric(20,8);not null;default:0"`
	RequestedAmountWei int       `gorm:"type:bigint;not null;default:0"`
	PaidAmountWei      int       `gorm:"type:bigint;not null;default:0"`
}
