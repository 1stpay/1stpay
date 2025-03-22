package kms

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/btcsuite/btcd/btcutil/base58"
)

type SolanaProvider struct{}

func (p SolanaProvider) Validate(address string) (bool, error) {
	return true, nil
}

func (p SolanaProvider) Create() (WalletData, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return WalletData{}, err
	}

	address := base58.Encode(publicKey)
	privateKeyStr := base58.Encode(privateKey)
	return WalletData{
		Address:    address,
		PrivateKey: privateKeyStr,
	}, nil
}
