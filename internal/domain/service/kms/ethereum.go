package kms

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumProvider struct{}

func (p EthereumProvider) Validate(address string) (bool, error) {
	return true, nil
}

func (p EthereumProvider) Create() (WalletData, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return WalletData{}, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return WalletData{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	privateKeyString := hex.EncodeToString(privateKeyBytes)
	privateKeyStringWithPrefix := "0x" + privateKeyString
	return WalletData{
		Address:    address,
		PrivateKey: privateKeyStringWithPrefix,
	}, nil
}
