package config

import (
	"github.com/1stpay/1stpay/internal/repository"
	"gorm.io/gorm"
)

type Repos struct {
	UserRepo       repository.UserRepositoryInterface
	MerchantRepo   repository.MerchantRepositoryInterface
	BlockchainRepo repository.BlockchainRepositoryInterface
	TokenRepo      repository.TokenRepositoryInterface
}

func NewRepositories(db *gorm.DB) *Repos {
	userRepo := repository.NewUserRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	blockchainRepo := repository.NewBlockchainRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	return &Repos{
		UserRepo:       userRepo,
		MerchantRepo:   merchantRepo,
		BlockchainRepo: blockchainRepo,
		TokenRepo:      tokenRepo,
	}
}
