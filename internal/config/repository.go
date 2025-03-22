package config

import (
	"github.com/1stpay/1stpay/internal/repository"
	"gorm.io/gorm"
)

type Repos struct {
	UserRepo           repository.UserRepository
	MerchantRepo       repository.MerchantRepository
	BlockchainRepo     repository.BlockchainRepository
	TokenRepo          repository.TokenRepository
	PaymentRepo        repository.PaymentRepository
	PaymentAddressRepo repository.PaymentAddressRepository
}

func NewRepositories(db *gorm.DB) *Repos {
	userRepo := repository.NewUserRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	blockchainRepo := repository.NewBlockchainRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	paymentAddressRepo := repository.NewPaymentAddressRepository(db)
	return &Repos{
		UserRepo:           userRepo,
		MerchantRepo:       merchantRepo,
		BlockchainRepo:     blockchainRepo,
		TokenRepo:          tokenRepo,
		PaymentRepo:        paymentRepo,
		PaymentAddressRepo: paymentAddressRepo,
	}
}
