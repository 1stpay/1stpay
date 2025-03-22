package blockchain_service

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	Client *ethclient.Client
}

func NewEthereumService(rpcURL string) (*EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &EthereumService{Client: client}, nil
}

func (s EthereumService) GetNativeBalance(address string) (*big.Int, error) {
	ctx := context.Background()
	account := common.HexToAddress(address)
	balanceWei, err := s.Client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balanceWei, nil
}

func (s EthereumService) GetTokenBalance(address, tokenAddress string) (*big.Int, error) {
	ctx := context.Background()

	account := common.HexToAddress(address)
	tokenAddr := common.HexToAddress(tokenAddress)

	methodID := common.Hex2Bytes("70a08231")
	paddedAddress := common.LeftPadBytes(account.Bytes(), 32)
	data := append(methodID, paddedAddress...)

	msg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("empty result from contract call")
	}

	balance := new(big.Int).SetBytes(result)
	return balance, nil
}
