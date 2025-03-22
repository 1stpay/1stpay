package infrastructure

import "math/big"

type BlockchainService interface {
	GetNativeBalance(address string) (*big.Int, error)
	GetTokenBalance(address, tokenAddress string) (*big.Int, error)
}
