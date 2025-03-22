package kms

import (
	"errors"

	"github.com/1stpay/1stpay/internal/domain/enum"
)

func GetProvider(chainType enum.NetworkType) (WalletProvider, error) {
	switch chainType {
	case enum.EVM:
		return EthereumProvider{}, nil
	case enum.TRON:
		return TronProvider{}, nil
	case enum.SOLANA:
		return SolanaProvider{}, nil
	default:
		return nil, errors.New("unsupported blockchain type")
	}
}
