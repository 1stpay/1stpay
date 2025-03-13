package config

import (
	"github.com/1stpay/1stpay/internal/domain/usecase"
)

type Usecases struct {
	AuthUsecase       usecase.AuthUsecaseInterface
	UserUsecase       usecase.UserUsecaseInterface
	MerchantUsecase   usecase.MerchantUsecaseInterface
	BlockchainUsecase usecase.BlockchainUsecaseInterface
	TokenUsecase      usecase.TokenUsecaseInterface
}

func NewUsecases(repos *Repos) *Usecases {
	userUsecase := usecase.NewUserUsecase(repos.UserRepo)
	authUsecase := usecase.NewAuthUsecase(repos.UserRepo)
	merchantUsecase := usecase.NewMerchantUsecase(repos.MerchantRepo)
	blockchainUsecase := usecase.NewBlockchainUsecase(repos.BlockchainRepo)
	tokenUsecase := usecase.NewTokenUsecase(repos.TokenRepo)
	return &Usecases{
		AuthUsecase:       authUsecase,
		UserUsecase:       userUsecase,
		MerchantUsecase:   merchantUsecase,
		BlockchainUsecase: blockchainUsecase,
		TokenUsecase:      tokenUsecase,
	}
}
