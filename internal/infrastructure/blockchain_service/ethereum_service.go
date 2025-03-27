package blockchain_service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	balanceOfMethodID    = common.Hex2Bytes("70a08231") // balanceOf(address)
	transferMethodID     = common.Hex2Bytes("a9059cbb") // transfer(address,uint256)
	errInvalidAddress    = errors.New("invalid address format")
	errZeroAmount        = errors.New("zero amount to transfer")
	errInsufficientFunds = errors.New("insufficient funds")
)

type EthereumService struct {
	ChainId *big.Int
	Client  *ethclient.Client
}

func NewEthereumService(rpcURL string, chainId int) (*EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	return &EthereumService{
		Client:  client,
		ChainId: big.NewInt(int64(chainId)),
	}, nil
}

func (s EthereumService) Close() {
	s.Client.Close()
}

// validateAddress validates an Ethereum address
func (s EthereumService) validateAddress(address string) (common.Address, error) {
	if !common.IsHexAddress(address) {
		return common.Address{}, errInvalidAddress
	}
	return common.HexToAddress(address), nil
}

// GetNativeBalance gets the balance of a native token for a specified address
func (s EthereumService) GetNativeBalance(ctx context.Context, address string) (*big.Int, error) {
	account, err := s.validateAddress(address)
	if err != nil {
		return nil, err
	}

	balanceWei, err := s.Client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balanceWei, nil
}

// GetTokenBalance gets the balance of a token for a specified address
func (s EthereumService) GetTokenBalance(ctx context.Context, address, tokenAddress string) (*big.Int, error) {
	account, err := s.validateAddress(address)
	if err != nil {
		return nil, err
	}

	tokenAddr, err := s.validateAddress(tokenAddress)
	if err != nil {
		return nil, err
	}

	paddedAddress := common.LeftPadBytes(account.Bytes(), 32)
	data := append(balanceOfMethodID, paddedAddress...)

	result, err := s.Client.CallContract(ctx, ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("empty result from contract call")
	}

	return new(big.Int).SetBytes(result), nil
}

// getSignerAndAddress returns the private key, address, and error
func (s *EthereumService) getSignerAndAddress(privateKey string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("public key format error")
	}

	return privKey, crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// TransferNative transfers a native token to a specified address
func (s *EthereumService) TransferNative(ctx context.Context, senderPrivateKey string, toAddress string, amount *big.Int) (common.Hash, error) {
	if amount.Sign() <= 0 {
		return common.Hash{}, errZeroAmount
	}

	recipient, err := s.validateAddress(toAddress)
	if err != nil {
		return common.Hash{}, err
	}

	privKey, fromAddress, err := s.getSignerAndAddress(senderPrivateKey)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := s.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %w", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &recipient,
		Value:    amount,
		Gas:      21000,
		GasPrice: gasPrice,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.ChainId), privKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := s.Client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, fmt.Errorf("transaction sending failed: %w", err)
	}

	return signedTx.Hash(), nil
}

// TransferNativeRemaining transfers the remaining balance of a native token to a specified address
func (s *EthereumService) TransferNativeRemaining(ctx context.Context, senderPrivateKey string, toAddress string) (common.Hash, error) {
	privKey, fromAddress, err := s.getSignerAndAddress(senderPrivateKey)
	if err != nil {
		return common.Hash{}, err
	}

	recipient, err := s.validateAddress(toAddress)
	if err != nil {
		return common.Hash{}, err
	}

	balance, err := s.GetNativeBalance(ctx, fromAddress.Hex())
	if err != nil {
		return common.Hash{}, fmt.Errorf("balance check failed: %w", err)
	}

	nonce, err := s.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("nonce check failed: %w", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("gas price check failed: %w", err)
	}

	fee := new(big.Int).Mul(big.NewInt(21000), gasPrice)
	if balance.Cmp(fee) <= 0 {
		return common.Hash{}, errInsufficientFunds
	}

	netAmount := new(big.Int).Sub(balance, fee)
	if netAmount.Sign() <= 0 {
		return common.Hash{}, errZeroAmount
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &recipient,
		Value:    netAmount,
		Gas:      21000,
		GasPrice: gasPrice,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.ChainId), privKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("signing failed: %w", err)
	}

	if err := s.Client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, fmt.Errorf("transaction sending failed: %w", err)
	}

	return signedTx.Hash(), nil
}

// TransferTokenRemaining transfers the remaining balance of a token to a specified address
func (s *EthereumService) TransferTokenRemaining(ctx context.Context, senderPrivateKey string, tokenAddress string, toAddress string) (common.Hash, error) {
	privKey, fromAddress, err := s.getSignerAndAddress(senderPrivateKey)
	if err != nil {
		return common.Hash{}, err
	}

	tokenAddr, err := s.validateAddress(tokenAddress)
	if err != nil {
		return common.Hash{}, err
	}

	recipient, err := s.validateAddress(toAddress)
	if err != nil {
		return common.Hash{}, err
	}

	tokenBalance, err := s.GetTokenBalance(ctx, fromAddress.Hex(), tokenAddr.Hex())
	if err != nil {
		return common.Hash{}, fmt.Errorf("token balance check failed: %w", err)
	}

	if tokenBalance.Sign() <= 0 {
		return common.Hash{}, errZeroAmount
	}

	paddedRecipient := common.LeftPadBytes(recipient.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(tokenBalance.Bytes(), 32)
	data := append(transferMethodID, append(paddedRecipient, paddedAmount...)...)

	nonce, err := s.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("nonce check failed: %w", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("gas price check failed: %w", err)
	}

	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &tokenAddr,
		GasPrice: gasPrice,
		Data:     data,
	}

	gasLimit, err := s.Client.EstimateGas(ctx, msg)
	if err != nil {
		gasLimit = 100000
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenAddr,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.ChainId), privKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("signing failed: %w", err)
	}

	if err := s.Client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, fmt.Errorf("transaction sending failed: %w", err)
	}

	return signedTx.Hash(), nil
}
