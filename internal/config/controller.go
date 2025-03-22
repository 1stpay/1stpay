package config

import (
	frontedController "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/controller"
)

type Controllers struct {
	MerchantAuthController       *controller.AuthController
	MerchantUserController       *controller.UserController
	MerchantMerchantController   *controller.MerchantController
	MerchantBlockchainController *controller.BlockchainController
	MerchantTokenController      *controller.TokenController
	MerchantPaymentController    *controller.PaymentController
	FrontendPaymentController    *frontedController.FrontendPaymentController
}

func NewControllers(usecases *Usecases) *Controllers {
	merchantAuthController := controller.NewAuthController(usecases.AuthUsecase)
	merchantUserController := controller.NewUserController(usecases.UserUsecase)
	merchantBlockchainController := controller.NewBlockchainController(usecases.BlockchainUsecase)
	merchantMerchantController := controller.NewMerchantController(usecases.MerchantUsecase, usecases.MerchantAPIKeyUsecase, usecases.UserUsecase)
	merchantTokenController := controller.NewTokenController(usecases.TokenUsecase)
	merchantPaymentController := controller.NewPaymentController(usecases.PaymentUsecase, usecases.MerchantUsecase, usecases.UserUsecase)
	frontedPaymentController := frontedController.NewPaymentController(usecases.PaymentUsecase)
	return &Controllers{
		MerchantAuthController:       merchantAuthController,
		MerchantUserController:       merchantUserController,
		MerchantMerchantController:   merchantMerchantController,
		MerchantBlockchainController: merchantBlockchainController,
		MerchantTokenController:      merchantTokenController,
		MerchantPaymentController:    merchantPaymentController,
		FrontendPaymentController:    frontedPaymentController,
	}
}
