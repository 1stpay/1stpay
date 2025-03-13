package config

import "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"

type Controllers struct {
	FrontendAuthController       *controller.AuthController
	FrontendUserController       *controller.UserController
	FrontendMerchantController   *controller.MerchantController
	FrontendBlockchainController *controller.BlockchainController
	FrontendTokenController      *controller.TokenController
}

func NewControllers(usecases *Usecases) *Controllers {
	frontendAuthController := controller.NewAuthController(usecases.AuthUsecase)
	frontendUserController := controller.NewUserController(usecases.UserUsecase)
	frontendBlockchainController := controller.NewBlockchainController(usecases.BlockchainUsecase)
	frontendMerchantController := controller.NewMerchantController(usecases.MerchantUsecase, usecases.UserUsecase)
	frontendTokenController := controller.NewTokenController(usecases.TokenUsecase)
	return &Controllers{
		FrontendAuthController:       frontendAuthController,
		FrontendUserController:       frontendUserController,
		FrontendMerchantController:   frontendMerchantController,
		FrontendBlockchainController: frontendBlockchainController,
		FrontendTokenController:      frontendTokenController,
	}
}
