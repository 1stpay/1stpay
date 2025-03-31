package blockchain_service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/1stpay/1stpay/pkg/tron"
	"github.com/ethereum/go-ethereum/common"
)

type TronService struct {
	client *tron.TronClient
}

func NewTronService(rpcURL string) (*TronService, error) {
	client := tron.NewTronClient(rpcURL)
	return &TronService{client: client}, nil
}

func (s TronService) GetNativeBalance(ctx context.Context, address string) (*big.Int, error) {
	return s.client.GetNativeBalance(ctx, address)
}

func (s TronService) GetTokenBalance(ctx context.Context, address, tokenAddress string) (*big.Int, error) {
	return s.client.GetTokenBalance(ctx, address, tokenAddress)
}

func (s TronService) TransferNative(ctx context.Context, senderPrivateKey string, toAddress string, amount *big.Int) (common.Hash, error) {
	return s.client.TransferNative(ctx, senderPrivateKey, toAddress, amount)
}

// TransferNativeRemaining transfers the entire native balance (minus fee) from the sender's wallet to the destination address.
func (s TronService) TransferNativeRemaining(ctx context.Context, senderPrivateKey string, toAddress string) (common.Hash, error) {
	fromAddress, err := tron.GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address from private key: %w", err)
	}
	balance, err := s.client.GetNativeBalance(ctx, fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address balance: %w", err)
	}
	return s.client.TransferNative(ctx, senderPrivateKey, toAddress, balance)
}

// TransferTokenRemaining transfers the entire balance of the specified TRC20 token from the sender's wallet to a destination address.
func (s TronService) TransferTokenRemaining(
	ctx context.Context,
	senderPrivateKey,
	tokenAddress,
	toAddress string,
) (common.Hash, error) {
	fromAddress, err := tron.GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address from private key: %w", err)
	}
	balance, err := s.client.GetTokenBalance(ctx, fromAddress, tokenAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address balance: %w", err)
	}
	return s.client.TransferToken(ctx, senderPrivateKey, tokenAddress, toAddress, balance)
}

// TransferTokenRemaining transfers the entire balance of the specified TRC20 token from the sender's wallet to a destination address.
func (s TronService) TransferToken(
	ctx context.Context,
	senderPrivateKey,
	tokenAddress,
	toAddress string,
	amount *big.Int,
) (common.Hash, error) {
	return s.client.TransferToken(ctx, senderPrivateKey, tokenAddress, toAddress, amount)
}
