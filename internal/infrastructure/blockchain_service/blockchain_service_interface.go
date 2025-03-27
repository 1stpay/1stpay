package blockchain_service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BlockchainService interface {
	// GetNativeBalance returns the native balance (in minimal units) for the given address.
	GetNativeBalance(ctx context.Context, address string) (*big.Int, error)
	// GetTokenBalance returns the token balance for the given wallet address and token contract address.
	GetTokenBalance(ctx context.Context, address, tokenAddress string) (*big.Int, error)
	// TransferETH sends a specified amount of native currency from the sender's wallet to a destination address.
	TransferNative(ctx context.Context, senderPrivateKey string, toAddress string, amount *big.Int) (common.Hash, error)
	// TransferETHRemaining transfers all native balance (minus fee) from the sender's wallet to the destination address.
	TransferNativeRemaining(ctx context.Context, senderPrivateKey string, toAddress string) (common.Hash, error)
	// TransferTokenRemaining transfers the entire balance of the specified ERC-20 token from the sender's wallet to a destination address.
	TransferTokenRemaining(ctx context.Context, senderPrivateKey string, tokenAddress string, toAddress string) (common.Hash, error)
}
