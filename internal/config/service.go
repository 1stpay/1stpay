package config

import (
	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
	"github.com/1stpay/1stpay/internal/infrastructure/price_service"
)

type Services struct {
	BlockchainService map[string]blockchain_service.BlockchainService
	PriceService      price_service.PriceService
}

func NewServices(repos *Repos, env *Env) *Services {
	blockchainService, err := blockchain_service.InitBlockchainServices(repos.BlockchainRepo)
	if err != nil {
		panic(err)
	}
	priceService := price_service.NewPriceService(env.PriceServiceKey)
	return &Services{
		BlockchainService: blockchainService,
		PriceService:      priceService,
	}
}
