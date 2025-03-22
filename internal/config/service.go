package config

import "github.com/1stpay/1stpay/internal/infrastructure"

type Services struct {
	BlockchainService map[string]infrastructure.BlockchainService
}

func NewServices(repos *Repos) *Services {
	blockchainService, err := infrastructure.InitBlockchainServices(repos.BlockchainRepo)
	if err != nil {
		panic(err)
	}
	return &Services{
		BlockchainService: blockchainService,
	}
}
