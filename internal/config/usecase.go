package config

import (
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"gorm.io/gorm"
)

type Usecases struct {
	AuthUsecase       usecase.AuthUsecaseInterface
	UserUsecase       usecase.UserUsecase
	MerchantUsecase   usecase.MerchantUsecase
	BlockchainUsecase usecase.BlockchainUsecaseInterface
	TokenUsecase      usecase.TokenUsecaseInterface
	PaymentUsecase    usecase.PaymentUsecaseInterface
}

func NewUsecases(db *gorm.DB, repos *Repos, services *Services) *Usecases {
	userUsecase := usecase.NewUserUsecase(repos.UserRepo)
	authUsecase := usecase.NewAuthUsecase(repos.UserRepo)
	merchantUsecase := usecase.NewMerchantUsecase(repos.MerchantRepo)
	blockchainUsecase := usecase.NewBlockchainUsecase(repos.BlockchainRepo)
	tokenUsecase := usecase.NewTokenUsecase(repos.TokenRepo)
	paymentUsecase := usecase.NewPaymentUsecase(db, repos.PaymentRepo, repos.PaymentAddressRepo, repos.MerchantRepo, services.PriceService)

	return &Usecases{
		AuthUsecase:       authUsecase,
		UserUsecase:       userUsecase,
		MerchantUsecase:   merchantUsecase,
		BlockchainUsecase: blockchainUsecase,
		TokenUsecase:      tokenUsecase,
		PaymentUsecase:    paymentUsecase,
	}
}
